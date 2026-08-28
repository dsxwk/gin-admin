package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

type CreateUserRolesTable struct{}

func (m *CreateUserRolesTable) ID() string {
	return "20251212_create_user_roles_table"
}

func (m *CreateUserRolesTable) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.UserRoles{})
}

func (m *CreateUserRolesTable) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.UserRoles{})
}
