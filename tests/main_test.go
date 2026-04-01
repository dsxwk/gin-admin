package tests

import (
	"gin/app/facade"
	"gin/bootstrap"
	"gin/pkg/foundation"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化app
	bootstrap.NewApp()

	// 初始化门面系统
	facade.Init()

	// 创建应用实例
	foundationApp := foundation.GetApp()
	facade.Register("app", foundationApp)

	// 启动应用(初始化所有providers)
	if err := foundationApp.Boot(); err != nil {
		panic("Failed to boot application: " + err.Error())
	}

	// 打印启动信息
	facade.Log.Info("Test environment initialized")

	// 运行测试
	code := m.Run()

	// 清理
	_ = foundationApp.Stop()

	os.Exit(code)
}
