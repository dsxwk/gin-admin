package make

import (
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"

	"github.com/samber/lo"
)

// MakeErrCode 错误码创建命令
type MakeErrCode struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeErrCode) Name() string {
	return "make:errcode"
}

// Description 命令描述
func (m *MakeErrCode) Description() string {
	return "错误码创建"
}

// Help 命令选项
func (m *MakeErrCode) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"模块名称, 如: login",
			true,
		},
		{
			base.Flag{
				Short: "p",
				Long:  "prefix",
			},
			"错误码前缀",
			false,
		},
	}
}

// Execute 执行命令
func (m *MakeErrCode) Execute(values map[string]string) {
	name := strings.TrimSuffix(strings.TrimSpace(values["file"]), ".go")
	if name == "" {
		m.ExitError("模块名称不能为空")
	}

	prefix := int64(0)
	if v := strings.TrimSpace(values["prefix"]); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			m.ExitError("错误码前缀必须为整数")
		}
		prefix = n
	} else {
		prefix = 0
	}

	name = filepath.ToSlash(name)
	baseName := filepath.Base(name)
	structName := lo.PascalCase(baseName) + "ErrCode"
	snakeName := lo.SnakeCase(baseName)
	file := filepath.Join("app", "errcode", name+"_errcode.go")
	m.generateFile(file, structName, snakeName, prefix)
}

// generateFile 生成错误码文件
func (m *MakeErrCode) generateFile(file, structName, snakeName string, prefix int64) {
	templateFile := m.GetTemplate("errcode")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("解析错误码模板失败: %s", err.Error())
		os.Exit(1)
	}

	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	data := struct {
		Package    string
		Name       string
		SnakeName  string
		Prefix     int64
		PrefixName string
	}{
		Package:    "errcode",
		Name:       structName,
		SnakeName:  snakeName,
		Prefix:     prefix,
		PrefixName: structName + "Prefix",
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("生成错误码文件失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("错误码文件: " + file + " 生成成功!")
}

func init() {
	cli.Register(&MakeErrCode{})
}
