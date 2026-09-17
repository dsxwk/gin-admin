package tests

import (
	"gin/app/facade"
	"gin/common/flag"
	_ "gin/common/imports"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 创建应用实例
	app := serviceprovider.NewApp()

	err := app.Boot()
	if err != nil {
		flag.Errorf("初始化应用失败: %v", err)
		os.Exit(1)
	}

	// 打印启动信息
	facade.Log().Info("Testing Initialized")
	errcode.SetLogger(facade.Log())

	// 运行测试
	code := m.Run()

	// 清理
	_ = app.Stop()

	os.Exit(code)
}
