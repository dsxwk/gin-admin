package seeds

import (
	"gin/app/model"

	"gorm.io/gorm"
)

type UserSeed struct{}

func (s *UserSeed) ID() string {
	return "20251212_user_seed"
}

func (s *UserSeed) Run(db *gorm.DB) error {
	users := []model.User{
		{
			Username: "admin",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Status:   1,
		},
		{
			Username: "test",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Status:   1,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	return nil
}
