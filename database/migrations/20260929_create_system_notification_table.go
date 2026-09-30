package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

// CreateSystemNotification20260929 系统通知表迁移
type CreateSystemNotification20260929 struct{}

// ID 迁移标识
func (m *CreateSystemNotification20260929) ID() string {
	return "20260929_create_system_notification_table"
}

// Migrate 创建通知表
func (m *CreateSystemNotification20260929) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.SystemNotification{})
}

// Rollback 删除通知表
func (m *CreateSystemNotification20260929) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.NotificationUsers{})
}
