package grpc

import (
	"bytes"
	"context"
	"fmt"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	protocGenGoVersion     = "v1.36.12"
	protocGenGoGrpcVersion = "v1.6.2"
)

// GrpcGen 生成grpc代码命令
type GrpcGen struct {
	base.BaseCommand
}

// Name 命令名称
func (g *GrpcGen) Name() string {
	return "grpc-gen"
}

// Description 命令描述
func (g *GrpcGen) Description() string {
	return "生成grpc代码"
}

// Help 命令选项
func (g *GrpcGen) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short:   "t",
				Long:    "type",
				Default: "all",
			},
			"生成类型,all/pb/grpc,默认all",
			false,
		},
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"proto文件路径,默认grpc/proto目录全部",
			false,
		},
		{
			base.Flag{
				Short: "d",
				Long:  "tool-dir",
			},
			"生成插件目录,默认.tools/bin",
			false,
		},
	}
}

// Execute 执行命令
func (g *GrpcGen) Execute(values map[string]string) {
	root := pkg.GetRootPath()
	genType := values["type"]
	if genType != "all" && genType != "pb" && genType != "grpc" {
		flag.Errorf("参数 --type 只能是 all/pb/grpc")
		return
	}

	files, err := g.protoFiles(root, values["file"])
	if err != nil {
		flag.Errorf("查找proto文件失败: %s", err.Error())
		return
	}
	if len(files) == 0 {
		flag.Errorf("未找到proto文件")
		return
	}

	toolDir := filepath.Join(root, ".tools", "bin")
	if dir := values["tool-dir"]; dir != "" {
		toolDir = g.absPath(root, dir)
	}
	if err = os.MkdirAll(toolDir, 0755); err != nil {
		flag.Errorf("创建插件目录失败: %s", err.Error())
		return
	}

	if genType == "all" || genType == "pb" {
		if err = g.ensureTool(toolDir, "protoc-gen-go", "google.golang.org/protobuf/cmd/protoc-gen-go", protocGenGoVersion); err != nil {
			flag.Errorf("安装protoc-gen-go失败: %s", err.Error())
			return
		}
	}
	if genType == "all" || genType == "grpc" {
		if err = g.ensureTool(toolDir, "protoc-gen-go-grpc", "google.golang.org/grpc/cmd/protoc-gen-go-grpc", protocGenGoGrpcVersion); err != nil {
			flag.Errorf("安装protoc-gen-go-grpc失败: %s", err.Error())
			return
		}
	}

	relNames := make([]string, len(files))
	for i, file := range files {
		rel, relErr := filepath.Rel(root, file)
		if relErr != nil {
			flag.Errorf("计算proto文件路径失败: %s", relErr.Error())
			return
		}
		relNames[i] = filepath.ToSlash(rel)
	}

	req, err := g.buildRequest(root, relNames)
	if err != nil {
		flag.Errorf("解析proto文件失败: %s", err.Error())
		return
	}

	if genType == "all" || genType == "pb" {
		if err = g.runPlugin(filepath.Join(toolDir, g.toolName("protoc-gen-go")), req, root); err != nil {
			flag.Errorf("生成pb代码失败: %s", err.Error())
			return
		}
	}
	if genType == "all" || genType == "grpc" {
		if err = g.runPlugin(filepath.Join(toolDir, g.toolName("protoc-gen-go-grpc")), req, root); err != nil {
			flag.Errorf("生成grpc代码失败: %s", err.Error())
			return
		}
	}

	flag.Successf("grpc代码生成完成")
}

// protoFiles 获取proto文件列表
func (g *GrpcGen) protoFiles(root, file string) ([]string, error) {
	if file == "" {
		var files []string
		seen := make(map[string]bool)
		patterns := []string{
			filepath.Join(root, "grpc", "proto", "*.proto"),
			filepath.Join(root, "grpc", "model", "*.proto"),
		}
		for _, pattern := range patterns {
			list, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			for _, item := range list {
				if !seen[item] {
					seen[item] = true
					files = append(files, item)
				}
			}
		}
		return files, nil
	}

	var files []string
	for item := range strings.SplitSeq(file, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		files = append(files, g.absPath(root, item))
	}
	return files, nil
}

// absPath 转为绝对路径
func (g *GrpcGen) absPath(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// buildRequest 构建代码生成请求
func (g *GrpcGen) buildRequest(root string, relNames []string) (*pluginpb.CodeGeneratorRequest, error) {
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{root},
		}),
	}

	files, err := compiler.Compile(context.Background(), relNames...)
	if err != nil {
		return nil, err
	}

	var all []protoreflect.FileDescriptor
	seen := make(map[string]bool)
	var collect func(protoreflect.FileDescriptor)
	collect = func(fd protoreflect.FileDescriptor) {
		if fd == nil || seen[fd.Path()] {
			return
		}
		seen[fd.Path()] = true
		for i := 0; i < fd.Imports().Len(); i++ {
			collect(fd.Imports().Get(i))
		}
		all = append(all, fd)
	}
	for _, file := range files {
		collect(file)
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: relNames,
		Parameter:      new("paths=source_relative"),
	}
	for _, fd := range all {
		req.ProtoFile = append(req.ProtoFile, protodesc.ToFileDescriptorProto(fd))
	}
	return req, nil
}

// ensureTool 安装生成插件
func (g *GrpcGen) ensureTool(toolDir, name, module, version string) error {
	path := filepath.Join(toolDir, g.toolName(name))
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	goBin, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("未找到go命令: %s", err.Error())
	}

	flag.Infof("安装%s %s", name, version)
	cmd := exec.Command(goBin, "install", module+"@"+version)
	cmd.Env = append(os.Environ(), "GOBIN="+toolDir)
	cmd.Dir = filepath.Dir(toolDir)
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return fmt.Errorf("%s: %s", runErr.Error(), string(output))
	}
	return nil
}

// runPlugin 运行生成插件
func (g *GrpcGen) runPlugin(tool string, req *pluginpb.CodeGeneratorRequest, outDir string) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	cmd := exec.Command(tool)
	cmd.Stdin = bytes.NewReader(data)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}

	resp := &pluginpb.CodeGeneratorResponse{}
	if err = proto.Unmarshal(out, resp); err != nil {
		return err
	}
	if resp.Error != nil && *resp.Error != "" {
		return fmt.Errorf("%s", *resp.Error)
	}

	for _, file := range resp.File {
		if file.GetName() == "" || file.GetContent() == "" {
			continue
		}
		path := filepath.Join(outDir, filepath.FromSlash(file.GetName()))
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err = os.WriteFile(path, []byte(file.GetContent()), 0644); err != nil {
			return err
		}
		flag.Successf("生成文件: " + path)
	}
	return nil
}

// toolName 获取插件文件名
func (g *GrpcGen) toolName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func init() {
	cli.Register(&GrpcGen{})
}
