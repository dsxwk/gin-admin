package database

import (
	"gin/app/model"

	"gorm.io/gorm"
)

// Migration Migration接口
type Migration interface {
	ID() string
	Migrate(db *gorm.DB) error
	Rollback(db *gorm.DB) error
}

// Migrations 迁移记录
type Migrations struct {
	ID        uint           `gorm:"primaryKey;comment:ID"`
	Migration string         `gorm:"size:191;not null;index:uniq_migration_type,unique;comment:迁移ID"`
	CreatedAt model.DateTime `gorm:"comment:创建时间"`
}

// TableOptions 表选项(可选)
func (Migrations) TableOptions() string {
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
}

// TableName 表名(可选)
func (Migrations) TableName() string {
	return "migrations"
}

// Seeder Seed接口
type Seeder interface {
	ID() string
	Run(db *gorm.DB) error
}
