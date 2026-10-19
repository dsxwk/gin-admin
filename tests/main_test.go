package tests

import (
	"gin/app/console"
	"gin/app/facade"
	"gin/tests/util"
	"strings"
	"testing"

	clibase "gin/pkg/cli/base"
)

func TestMain(m *testing.M) {
	util.Run(m)
}

// TestCommandCheckDirAndFileRejectsTraversal 测试文件路径穿越校验
func TestCommandCheckDirAndFileRejectsTraversal(t *testing.T) {
	command := &clibase.Command{}
	if file := command.CheckDirAndFile("../secret.go"); file != nil {
		_ = file.Close()
		t.Fatal("非法路径不应创建文件")
	}
}

// TestCommandFormatArgs 测试参数格式稳定
func TestCommandFormatArgs(t *testing.T) {
	command := &clibase.Command{}
	got := command.FormatArgs(map[string]string{
		"name": "user",
		"desc": "用户",
		"skip": "",
	})
	if !strings.HasPrefix(got, "--desc=用户 --name=user") {
		t.Fatalf("参数顺序错误:%s", got)
	}
	if strings.Contains(got, "skip") {
		t.Fatalf("空参数不应输出:%s", got)
	}
}

// TestCommandParseFlagsReturnsError 测试参数解析错误返回
func TestCommandParseFlagsReturnsError(t *testing.T) {
	command := &clibase.Command{}
	_, err := command.ParseFlags("test", []string{"--unknown"}, nil)
	if err == nil {
		t.Fatal("未知参数应返回错误")
	}
}

// TestCommandConfig 测试命令配置
func TestCommandConfig(t *testing.T) {
	command := clibase.NewCommand(clibase.Config{
		Name:        "demo:test",
		Description: "演示控制台命令",
		Options: []clibase.CommandOption{
			{
				Flag: clibase.Flag{Short: "a", Long: "args"},
				Desc: "演示参数",
			},
		},
	})

	if command.Name() != "demo:test" {
		t.Fatalf("命令名称错误:%s", command.Name())
	}
	if command.Description() != "演示控制台命令" {
		t.Fatalf("命令描述错误:%s", command.Description())
	}
	if len(command.Help()) != 1 || command.Help()[0].Desc != "演示参数" {
		t.Fatal("命令帮助配置错误")
	}
}

// TestConsoleCommandConfig 测试控制台命令配置
func TestConsoleCommandConfig(t *testing.T) {
	command := console.NewDemoCommand()
	if command.Name() != "demo:test" {
		t.Fatalf("命令名称错误:%s", command.Name())
	}
	if command.Description() != "演示控制台命令" {
		t.Fatalf("命令描述错误:%s", command.Description())
	}
	if len(command.Help()) != 1 {
		t.Fatal("命令帮助配置错误")
	}
}

// TestAgentFacade 测试AI Agent门面
func TestAgentFacade(t *testing.T) {
	manager := facade.Agent()
	if manager == nil {
		t.Fatal("AI Agent服务未注册")
	}
	providerName, modelName, ok := manager.Provider()
	if !ok {
		t.Fatal("AI Agent默认提供商未配置")
	}
	if providerName == "" || modelName == "" {
		t.Fatal("AI Agent默认提供商信息不完整")
	}
}
