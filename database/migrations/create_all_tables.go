package migrations

import (
	"gin/database"
	"gin/database/seeds"
)

var (
	migration []database.Migration // 注册迁移
	seed      []database.Seeder    // 注册 seed
)

func init() {
	// 注册迁移
	migration = []database.Migration{
		&CreateUserTable{},
		&CreateUserRolesTable{},
		&CreateAgentSessionTable{},
		&CreateAgentMessageTable{},
		// ...
	}

	// 注册 seed
	seed = []database.Seeder{
		&seeds.UserSeed{},
		&seeds.UserRolesSeed{},
		// ...
	}
}

func AllMigrations() []database.Migration {
	return migration
}

func AllSeeds() []database.Seeder {
	return seed
}
