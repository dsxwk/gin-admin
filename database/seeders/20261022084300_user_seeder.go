package seeders

import (
	"context"

	"gin/app/model"

	"gorm.io/gorm"
)

type UserSeeder struct{}

// Name Seeder名称
func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

// Run 写入默认用户
func (s *UserSeeder) Run(ctx context.Context, db *gorm.DB) error {
	users := []model.User{
		{
			Avatar:   "",
			Username: "admin",
			FullName: "超级管理员",
			Nickname: "大师兄",
			Email:    "dsx.emil@qq.com",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Gender:   1,
			Age:      31,
			Status:   1,
		},
		{
			Avatar:   "",
			Username: "test",
			FullName: "测试管理员",
			Nickname: "dsx",
			Email:    "dsx1.emil@qq.com",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Gender:   1,
			Age:      31,
			Status:   1,
		},
	}

	for index := range users {
		if err := db.WithContext(ctx).
			Where("username = ?", users[index].Username).
			FirstOrCreate(&users[index]).
			Error; err != nil {
			return err
		}
	}
	return nil
}
