package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

type CreateUserTable struct{}

func (m *CreateUserTable) ID() string {
	return "20251212_create_user_table"
}

func (m *CreateUserTable) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}

func (m *CreateUserTable) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.User{})
}
