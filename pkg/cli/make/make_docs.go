package make

import (
	"fmt"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os/exec"
)

type MakeDocs struct {
	base.Command
}

// NewMakeDocs 创建命令
func NewMakeDocs() *MakeDocs {
	return &MakeDocs{
		Command: base.NewCommand(base.Config{
			Name:        "make:docs",
			Description: "生成Swagger文档",
		}),
	}
}

func (m *MakeDocs) Handle(_ map[string]string) error {
	terminal.Infof("开始生成Swagger文档...")

	// 排除远程调用服务目录,避免同名模型导致文档解析失败
	cmd := exec.Command("swag", "init", "-g", "main.go", "--exclude", "grpc")
	output, err := cmd.CombinedOutput()
	if err != nil {
		terminal.Errorf("Swagger文档生成失败: %v\n%s", err, string(output))
		return nil

	}

	fmt.Println(string(output))
	terminal.Successf("Swagger文档生成成功")
	return nil
}
