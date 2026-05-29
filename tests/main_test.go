package tests

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/common/response"
	"gin/pkg/serviceprovider"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化门面系统
	facade.Init()
	// 创建应用实例
	app := serviceprovider.GetApp()
	// 注册应用到门面
	facade.Register("app", app)

	// 打印启动信息
	facade.Log().Info("Test environment initialized")
	err := app.Boot()
	if err != nil {
		flag.Errorf("初始化应用失败: %v", err)
		os.Exit(1)
	}

	response.SetLogger(facade.Log())

	// 运行测试
	code := m.Run()

	// 清理
	_ = app.Stop()

	os.Exit(code)
}
