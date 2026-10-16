package mcp

import (
	"bytes"
	"context"
	"fmt"
	"gin/app/console"
	"gin/pkg/cli"
	"gin/pkg/cli/registry"
	"gin/pkg/serviceprovider/mcp"
	"io"
	"os"
	"slices"
	"strings"
)

// CliExecArgs CLI命令执行参数
type CliExecArgs struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

// NewCliExecTool 创建CLI命令执行工具
func NewCliExecTool() *mcp.ToolAdapter[CliExecArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name: "execute_cli",
			Description: `执行CLI命令,可用命令:
				make:controller(创建控制器/file=路径,desc=描述),
				make:model(创建模型/name=名称,table=表名),
				make:service(创建服务/name=名称),
				make:request(创建验证器/name=名称),
				make:middleware(创建中间件/name=名称),
				make:router(创建路由/file=路径),
				make:docs(生成Swagger文档),
				permission:sync(同步用户权限到Redis),
				route:list(查看路由列表),
				job:list(查看Job列表)`,
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"command": {Type: "string", Description: "CLI命令名称,如make:controller"},
					"args":    {Type: "string", Description: "命令参数,格式: name=User&file=v1/user&desc=用户管理,多个参数用&分隔"},
				},
				Required: []string{"command"},
			},
		},
		handleCliExec,
	)
}

// handleCliExec 处理CLI命令执行工具
func handleCliExec(_ context.Context, _ int64, args *CliExecArgs) (any, error) {
	if args.Command == "" {
		return map[string]any{"error": "请指定命令名称,如make:controller"}, nil
	}

	commands := cli.NewRegistry(slices.Concat(console.All(), registry.All())...)
	cmd, exists := commands.Get(args.Command)
	if !exists {
		available := listAvailableCliCommands()
		return map[string]any{"error": fmt.Sprintf("命令 %s 不存在。可用命令: %s", args.Command, available)}, nil
	}

	values := make(map[string]string)
	if args.Args != "" {
		for pair := range strings.SplitSeq(args.Args, "&") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				values[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			} else if len(kv) == 1 {
				values[strings.TrimSpace(kv[0])] = ""
			}
		}
	}

	if name, ok := values["name"]; ok && values["file"] == "" {
		values["file"] = name
	}

	output := captureCliOutput(func() {
		cmd.Execute(values)
	})

	resultFile := buildCliResultPath(args.Command, values)
	result := map[string]any{
		"command": args.Command,
		"args":    values,
		"output":  output,
	}
	if resultFile != "" {
		result["created"] = resultFile
	}
	result["message"] = fmt.Sprintf("命令 %s 执行完成", args.Command)

	return result, nil
}

// captureCliOutput 捕获CLI标准输出
func captureCliOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = old

	return <-done
}

// listAvailableCliCommands 获取可用CLI命令
func listAvailableCliCommands() string {
	names := []string{
		"make:controller", "make:model", "make:service", "make:request",
		"make:middleware", "make:router", "make:docs",
		"permission:sync", "permission:grant-admin",
		"route:list", "job:list",
	}
	return strings.Join(names, ", ")
}

// buildCliResultPath 构建CLI生成文件路径
func buildCliResultPath(cmdName string, values map[string]string) string {
	file, _ := values["file"]
	name, _ := values["name"]
	if file == "" {
		file = name
	}
	if file == "" {
		return ""
	}

	switch cmdName {
	case "make:controller":
		return "app/controller/" + file + ".go"
	case "make:model":
		return "app/model/" + file + ".go"
	case "make:service":
		return "app/service/" + file + ".go"
	case "make:request":
		return "app/request/" + file + ".go"
	case "make:middleware":
		return "app/middleware/" + file + ".go"
	case "make:router":
		return "router/" + file + ".go"
	default:
		return ""
	}
}
