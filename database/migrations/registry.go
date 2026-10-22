package migrations

import "gin/pkg/migration"

// All 获取全部迁移
func All() []migration.Migration {
	return []migration.Migration{
		&CreateUsersTable{},
	}
}
