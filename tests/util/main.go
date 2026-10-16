package util

import (
	"gin/app/facade"
	"gin/app/provider"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider"
	"gin/pkg/terminal"
	"os"
	"testing"
)

// Run 初始化测试应用并执行测试
func Run(m *testing.M) {
	app := serviceprovider.NewApp(provider.All()...)
	if err := app.Boot(); err != nil {
		terminal.Errorf("初始化应用失败: %v", err)
		os.Exit(1)
	}

	facade.Log().Info("Testing Initialized")
	errcode.SetLogger(facade.Log())

	code := m.Run()
	_ = app.Stop()
	os.Exit(code)
}
