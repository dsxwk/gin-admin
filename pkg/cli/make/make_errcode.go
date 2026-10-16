package make

import (
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gin/pkg/cli/base"
	"gin/pkg/terminal"

	"github.com/samber/lo"
)

// MakeErrCode 错误码创建命令
type MakeErrCode struct {
	base.Command
}

// NewMakeErrCode 创建命令
func NewMakeErrCode() *MakeErrCode {
	return &MakeErrCode{
		Command: base.NewCommand(base.Config{
			Name:        "make:errcode",
			Description: "错误码创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "模块名称, 如: login",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "p", Long: "prefix"},
					Desc:     "错误码前缀",
					Required: false,
				},
			},
		}),
	}
}

// Handle 执行命令
func (m *MakeErrCode) Handle(values map[string]string) error {
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
	return nil
}

// generateFile 生成错误码文件
func (m *MakeErrCode) generateFile(file, structName, snakeName string, prefix int64) {
	templateFile := m.GetTemplate("errcode")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("解析错误码模板失败: %s", err.Error())
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
		terminal.Errorf("生成错误码文件失败: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("错误码文件: " + file + " 生成成功!")
}
