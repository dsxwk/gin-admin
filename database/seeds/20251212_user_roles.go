package seeds

import (
	"gin/app/model"

	"gorm.io/gorm"
)

type UserRolesSeed struct{}

func (s *UserRolesSeed) ID() string {
	return "20251212_user_roles_seed"
}

func (s *UserRolesSeed) Run(db *gorm.DB) error {
	userRoles := []model.UserRoles{
		{
			UserID: 1,
			RoleID: 1,
			Name:   "",
		},
		{
			UserID: 2,
			RoleID: 1,
			Name:   "",
		},
	}

	if err := db.Create(&userRoles).Error; err != nil {
		return err
	}

	return nil
}
