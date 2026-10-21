package tests

import (
	"context"
	"gin/app/console"
	"gin/app/facade"
	"gin/pkg/cli/base"
	"gin/tests/util"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	util.Run(m)
}

// TestCommandCheckDirAndFileRejectsTraversal 测试文件路径穿越校验
func TestCommandCheckDirAndFileRejectsTraversal(t *testing.T) {
	command := &base.Command{}
	if file := command.CheckDirAndFile("../secret.go"); file != nil {
		_ = file.Close()
		t.Fatal("非法路径不应创建文件")
	}
}

// TestCommandFormatArgs 测试参数格式稳定
func TestCommandFormatArgs(t *testing.T) {
	command := &base.Command{}
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
	command := &base.Command{}
	_, err := command.ParseFlags("test", []string{"--unknown"}, nil)
	if err == nil {
		t.Fatal("未知参数应返回错误")
	}
}

// TestCommandConfig 测试命令配置
func TestCommandConfig(t *testing.T) {
	command := base.NewCommand(base.Config{
		Name:        "demo:test",
		Description: "演示控制台命令",
		Options: []base.CommandOption{
			{
				Flag: base.Flag{Short: "a", Long: "args"},
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
	if command.Help()[0].Flag.Long != "args" {
		t.Fatal("普通命令不应自动添加cron参数")
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
	if len(command.Help()) != 2 {
		t.Fatal("命令帮助配置错误")
	}
	if command.Help()[1].Flag.Long != "cron" {
		t.Fatal("控制台命令缺少cron参数")
	}
}

// TestConsoleCommandsSupportCron 测试所有控制台命令支持cron
func TestConsoleCommandsSupportCron(t *testing.T) {
	for _, command := range console.All() {
		supported := false
		for _, option := range command.Help() {
			if option.Flag.Long == "cron" {
				supported = true
				break
			}
		}
		if !supported {
			t.Fatalf("命令%s不支持cron", command.Name())
		}
	}
}

// TestDemoCommandCron 测试演示命令cron表达式
func TestDemoCommandCron(t *testing.T) {
	command := console.NewDemoCommand()
	if err := base.ExecuteWithCron(command, map[string]string{"args": "test"}); err != nil {
		t.Fatalf("空cron表达式执行失败:%v", err)
	}
	if err := base.ExecuteWithCron(command, map[string]string{"args": "test", "cron": "invalid"}); err == nil {
		t.Fatal("非法cron表达式应返回错误")
	}
}

// TestCommandCronAliases 测试cron参数长短写法
func TestCommandCronAliases(t *testing.T) {
	command := console.NewDemoCommand()
	for _, cronFlag := range []string{"--cron", "-cron"} {
		values, err := command.ParseFlags("demo:test", []string{"--args=test", cronFlag, "*/5 * * * * *"}, command.Help())
		if err != nil {
			t.Fatalf("cron参数%s解析失败:%v", cronFlag, err)
		}
		if values["cron"] != "*/5 * * * * *" {
			t.Fatalf("cron参数%s解析结果错误:%s", cronFlag, values["cron"])
		}
	}
}

// TestDemoCommandCronMetadata 测试演示命令自动调度配置
func TestDemoCommandCronMetadata(t *testing.T) {
	command := console.NewDemoCommand()
	if command.Cron() != "*/5 * * * * *" {
		t.Fatalf("cron表达式错误:%s", command.Cron())
	}
	if command.CronValues()["args"] != "arg1" {
		t.Fatal("cron参数配置错误")
	}
}

// TestCronManagerCollectsCommands 测试定时任务自动收集
func TestCronManagerCollectsCommands(t *testing.T) {
	manager, err := base.NewCronManager(console.All()...)
	if err != nil {
		t.Fatalf("创建cron管理器失败:%v", err)
	}
	if manager.Len() == 0 {
		t.Fatal("未收集到定时命令")
	}

	invalid := &cronInvalidCommand{
		Command: base.NewCommand(base.Config{
			Name: "invalid:cron",
			Cron: "invalid",
		}),
	}
	if _, err = base.NewCronManager(invalid); err == nil {
		t.Fatal("非法cron表达式应返回错误")
	}
}

// TestCronManagerAutoExecute 测试定时任务自动执行
func TestCronManagerAutoExecute(t *testing.T) {
	command := &cronCountCommand{
		Command: base.NewCommand(base.Config{
			Name: "count:cron",
			Cron: "@every 1s",
		}),
	}
	manager, err := base.NewCronManager(command)
	if err != nil {
		t.Fatalf("创建cron管理器失败:%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- manager.Run(ctx)
	}()

	deadline := time.Now().Add(2500 * time.Millisecond)
	for command.count.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	cancel()
	if err = <-done; err != nil {
		t.Fatalf("运行cron管理器失败:%v", err)
	}
	if err = manager.Stop(); err != nil {
		t.Fatalf("停止cron管理器失败:%v", err)
	}
	if command.count.Load() == 0 {
		t.Fatal("定时任务未自动执行")
	}
}

type cronInvalidCommand struct {
	base.Command
}

// Handle 执行测试命令
func (c *cronInvalidCommand) Handle(_ map[string]string) error {
	return nil
}

type cronCountCommand struct {
	base.Command
	count atomic.Int32
}

// Handle 执行计数命令
func (c *cronCountCommand) Handle(_ map[string]string) error {
	c.count.Add(1)
	return nil
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
