package make

import (
	"fmt"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"os/exec"
)

type MakeDocs struct {
	base.BaseCommand
}

func (m *MakeDocs) Name() string {
	return "make:docs"
}

func (m *MakeDocs) Description() string {
	return "生成Swagger文档"
}

func (m *MakeDocs) Execute(values map[string]string) {
	flag.Infof("开始生成Swagger文档...")

	cmd := exec.Command("swag", "init", "-g", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		flag.Errorf("Swagger文档生成失败: %v\n%s", err, string(output))
		return
	}

	fmt.Println(string(output))
	flag.Successf("Swagger文档生成成功")
}

func init() {
	cli.Register(&MakeDocs{})
}
