package orm

import (
	"testing"

	"gin/config"
	serviceorm "gin/pkg/serviceprovider/orm"

	"github.com/stretchr/testify/require"
)

// TestManagerConnectUnsupported 测试不支持驱动返回错误
func TestManagerConnectUnsupported(t *testing.T) {
	manager := serviceorm.NewManager(&config.Config{})
	db, err := manager.Connect("unsupported")
	require.Error(t, err)
	require.Nil(t, db)
}

// TestManagerConnectionUnsupported 测试兼容连接方法
func TestManagerConnectionUnsupported(t *testing.T) {
	manager := serviceorm.NewManager(&config.Config{})
	require.Nil(t, manager.Connection("unsupported"))
}
