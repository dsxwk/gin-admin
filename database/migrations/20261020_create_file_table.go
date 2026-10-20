package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

// CreateFile20261020 文件表和分组表迁移
type CreateFile20261020 struct{}

// ID 迁移标识
func (m *CreateFile20261020) ID() string {
	return "20261020_create_file_table"
}

// Migrate 创建文件表和分组表
func (m *CreateFile20261020) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.FileGroup{}, &model.File{})
}

// Rollback 删除文件表和分组表
func (m *CreateFile20261020) Rollback(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&model.File{}); err != nil {
		return err
	}
	return db.Migrator().DropTable(&model.FileGroup{})
}
