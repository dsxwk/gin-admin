package migration

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Migration 数据库迁移接口
type Migration interface {
	Name() string
	Up(ctx context.Context, db *gorm.DB) error
	Down(ctx context.Context, db *gorm.DB) error
}

// Status 迁移状态
type Status struct {
	Name      string
	Batch     int
	Ran       bool
	CreatedAt *time.Time
}

// RunError 迁移执行错误
type RunError struct {
	Name string
	Err  error
}

// Error 返回迁移错误信息
func (e *RunError) Error() string {
	return e.Name + " 执行失败: " + e.Err.Error()
}

// Unwrap 获取原始错误
func (e *RunError) Unwrap() error {
	return e.Err
}
