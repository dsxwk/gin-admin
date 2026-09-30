package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

type CreateNotificationUsers20260929 struct{}

func (m *CreateNotificationUsers20260929) ID() string {
	return "20260929_create_notification_users_table"
}

func (m *CreateNotificationUsers20260929) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.NotificationUsers{})
}

func (m *CreateNotificationUsers20260929) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.NotificationUsers{})
}
