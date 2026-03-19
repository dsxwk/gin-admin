package tests

import (
	"gin/bootstrap"
	"gin/config"
	"gin/pkg/container"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化配置
	_ = config.NewConfig()

	// 初始化container
	c := container.NewContainer()

	// 初始化app
	_ = bootstrap.NewApp(c)

	// 运行测试
	code := m.Run()

	os.Exit(code)
}
