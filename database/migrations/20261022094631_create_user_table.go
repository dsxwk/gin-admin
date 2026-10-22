package migrations

import (
	"context"

	"gin/pkg/migration"

	"gorm.io/gorm"
)

type CreateUsersTable struct{}

// Name 迁移名称
func (m *CreateUsersTable) Name() string {
	return "20261022094631_create_user_table"
}

// Up 执行迁移
func (m *CreateUsersTable) Up(ctx context.Context, db *gorm.DB) error {
	return migration.CreateTable(db.WithContext(ctx), "user", func(t *migration.Table) {
		t.Id()
		t.String("avatar", 255).NotNull()
		t.String("username", 50).NotNull().Unique()
		t.String("full_name", 50).NotNull()
		t.String("email", 100).NotNull()
		t.String("password", 255).NotNull()
		t.String("nickname", 50).NotNull()
		t.Int("gender").Default(0).NotNull()
		t.Int("age").Default(0).NotNull()
		t.Int("status").Default(1).NotNull()
		t.Timestamps()
		t.SoftDeletes()
	})
}

// Down 回滚迁移
func (m *CreateUsersTable) Down(ctx context.Context, db *gorm.DB) error {
	return migration.DropTable(db.WithContext(ctx), "user")
}
