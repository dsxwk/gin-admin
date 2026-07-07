package service

import (
	"errors"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"github.com/samber/lo"
	"time"
)

type UserService struct {
	base.BaseService
}

// List 列表
func (s *UserService) List(req request.User) (pageData request.PageData, err error) {
	var (
		m  []model.User
		db = s.DB(&model.User{})
	)

	// 搜索
	db = s.Search(db, m, req.Search)

	err = db.Model(&m).Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	db = db.Model(&m).Preload("UserRoles")

	if req.NotPage {
		err = db.Order("id DESC").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m
	} else {
		pageData.Page = req.Page
		pageData.PageSize = req.PageSize
		offset, limit := request.Pagination(req.Page, req.PageSize)

		err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m
	}

	return pageData, nil
}

// Create 创建
func (s *UserService) Create(req request.User) (m model.User, err error) {
	var (
		count     int64
		userRoles []model.UserRoles
		db        = s.DB(&model.User{})
	)

	// 校验用户名是否重复
	err = db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error
	if err != nil {
		return m, err
	}
	if count > 0 {
		return m, errors.New("用户名已存在")
	}

	m = model.User{
		Username: req.Username,
		FullName: req.FullName,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Age:      req.Age,
	}

	// 处理密码
	m.Password = pkg.BcryptHash(req.Password)

	tx := db.Begin()

	err = tx.Model(&model.User{}).Create(&m).Error
	if err != nil {
		tx.Rollback()
		return m, err
	}

	if req.UserRoles != nil {
		for _, v := range req.UserRoles {
			userRoles = append(userRoles, model.UserRoles{
				UserID: m.ID,
				RoleID: v.RoleId,
				Name:   v.Name,
			})
		}
		err = tx.Model(&model.UserRoles{}).Create(&userRoles).Error
		if err != nil {
			tx.Rollback()
			return m, err
		}
	}

	tx.Commit()

	return m, nil
}

// Update 更新
func (s *UserService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		count int64
		db    = s.DB(&model.User{})
	)

	userRolesData, ok := data["userRoles"].([]interface{})
	if !ok {
		userRolesData = []interface{}{}
	}

	// 校验用户名是否重复
	err = db.Model(&model.User{}).Where("username = ? AND id <> ?", data["username"], id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	if lo.HasKey(data, "password") && data["password"] != "" {
		data["password"] = pkg.BcryptHash(data["password"].(string))
	}
	rows := model.FilterFields(db, model.User{}, data)
	rows["updated_at"] = time.DateTime

	tx := s.DB(&model.User{}).Begin()

	err = tx.Model(&model.User{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新用户角色关联
	if len(userRolesData) > 0 {
		// 删除旧关联
		err = tx.Model(&model.UserRoles{}).
			Where("user_id = ?", id).
			Delete(&model.UserRoles{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		var newUserRoles []model.UserRoles
		for _, item := range userRolesData {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newUserRoles = append(newUserRoles, model.UserRoles{
				UserID: id,
				RoleID: int64(roleMap["roleId"].(float64)),
				Name:   roleMap["name"].(string),
			})
		}

		if len(newUserRoles) > 0 {
			err = tx.Model(&model.UserRoles{}).Create(&newUserRoles).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// Detail 详情
func (s *UserService) Detail(id int64) (m model.User, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).
		Preload("UserRoles").
		First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Delete 删除
func (s *UserService) Delete(id int64) (err error) {
	var (
		m  model.User
		db = s.DB(&m)
	)

	err = db.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}

	return nil
}
