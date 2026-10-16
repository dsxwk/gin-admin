package orm

import (
	"strings"
	"testing"

	baselib "gin/common/base"
	"gin/config"
	serviceorm "gin/pkg/serviceprovider/orm"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

type searchTestModel struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
}

// TableName 测试表名
func (searchTestModel) TableName() string {
	return "search_test_model"
}

// TestBaseServiceSearchReturnsBuildError 测试搜索条件错误向上返回
func TestBaseServiceSearchReturnsBuildError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&searchTestModel{}))

	result := (&baselib.Service{}).Search(db, &searchTestModel{}, map[string]any{"missing": "value"})
	require.Error(t, result.Error)
	require.True(t, strings.Contains(result.Error.Error(), "字段"))

	var rows []searchTestModel
	require.Error(t, result.Model(&searchTestModel{}).Find(&rows).Error)
}
