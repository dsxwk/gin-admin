package seeders

import "gin/pkg/migration"

// All 获取全部Seeder
func All() []migration.Seeder {
	return []migration.Seeder{
		&UserSeeder{},
	}
}
