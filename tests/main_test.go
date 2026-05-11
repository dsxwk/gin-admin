package tests

import (
	"gin/app/facade"
	"gin/pkg/foundation"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化门面系统
	facade.Init()
	// 创建应用实例
	app := foundation.GetApp()
	// 注册应用到门面
	facade.Register("app", app)

	// 打印启动信息
	facade.Log().Info("Test environment initialized")

	// 运行测试
	code := m.Run()

	// 清理
	_ = app.Stop()

	os.Exit(code)
}
