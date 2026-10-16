package migrations

import (
	"gin/database"
	"gin/database/seeds"
)

// Migrations 迁移列表
func Migrations() []database.Migration {
	return []database.Migration{
		&CreateUserTable{},
		&CreateUserRolesTable{},
		&CreateAgentSessionTable{},
		&CreateAgentMessageTable{},
	}
}

// Seeds 数据填充列表
func Seeds() []database.Seeder {
	return []database.Seeder{
		&seeds.UserSeed{},
		&seeds.UserRolesSeed{},
	}
}
