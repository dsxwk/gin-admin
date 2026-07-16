package service

import (
	"errors"
	"gin/app/enum"
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
	db = s.Search(db, m, req.Search).Model(&m).Preload("UserRoles")

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

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
	rows["updated_at"] = time.Now()

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

// BatchDelete 批量删除
func (s *UserService) BatchDelete(ids []int64) (err error) {
	var (
		m  model.User
		db = s.DB(&m)
	)

	err = db.Model(&m).Delete(&m, ids).Error
	if err != nil {
		return err
	}

	return nil
}

// Import 批量导入用户
func (s *UserService) Import(req request.UserImport) (request.UserImport, error) {
	db := s.DB(&model.User{})

	// 收集所有用户名和邮箱
	usernames := make([]string, 0, len(req.Data))
	emails := make([]string, 0, len(req.Data))
	for _, item := range req.Data {
		usernames = append(usernames, item.Username)
		emails = append(emails, item.Email)
	}

	// 批量查询数据库中已存在的用户名和邮箱
	existingUsernames := make(map[string]bool)
	existingEmails := make(map[string]bool)

	var existedUsers []model.User
	if len(usernames) > 0 || len(emails) > 0 {
		db.Model(&model.User{}).Where("username IN ? OR email IN ?", usernames, emails).Find(&existedUsers)
	}
	for _, u := range existedUsers {
		existingUsernames[u.Username] = true
		existingEmails[u.Email] = true
	}

	// 逐行检查并收集错误
	seenUsername := make(map[string]int)
	seenEmail := make(map[string]int)
	var batchUsers []model.User
	for i, item := range req.Data {
		row := i + 1

		// 校验导入数据内用户名重复
		if firstRow, ok := seenUsername[item.Username]; ok {
			return req, errors.New("第 " + pkg.IntToString(row) + " 行,用户名 " + item.Username + " 重复(与第" + pkg.IntToString(firstRow) + "行)")
		}
		// 校验导入数据内邮箱重复
		if firstRow, ok := seenEmail[item.Email]; ok {
			return req, errors.New("第 " + pkg.IntToString(row) + " 行,邮箱 " + item.Email + " 重复(与第" + pkg.IntToString(firstRow) + "行)")
		}
		if existingUsernames[item.Username] {
			return req, errors.New("第 " + pkg.IntToString(row) + " 行,用户名 " + item.Username + " 已存在")
		}
		if existingEmails[item.Email] {
			return req, errors.New("第 " + pkg.IntToString(row) + " 行,邮箱 " + item.Email + " 已存在")
		}

		// 记录已处理的行号,供后续行检查内部重复
		seenUsername[item.Username] = row
		seenEmail[item.Email] = row

		batchUsers = append(batchUsers, model.User{
			Username: item.Username,
			Password: pkg.BcryptHash(item.Password),
			FullName: item.FullName,
			Nickname: item.Nickname,
			Email:    item.Email,
			Gender:   item.Gender,
			Age:      item.Age,
			Status:   item.Status,
		})
	}
	tx := db.Begin()
	if len(batchUsers) > 0 {
		if err := db.Create(&batchUsers).Error; err != nil {
			tx.Rollback()
			return req, err
		}

		var (
			importModel model.ImportRecords
			importEnum  enum.ImportRecordsEnum
		)
		data := model.JsonValue{
			Data: req.Data,
		}
		err := tx.Model(&importModel).Create(&model.ImportRecords{
			Type: enum.ImportRecordsTypeUser,
			Name: importEnum.Type().Desc(enum.ImportRecordsTypeUser),
			Data: &data,
		}).Error
		if err != nil {
			tx.Rollback()
			return req, err
		}
	}

	tx.Commit()

	return req, nil
}

// Password 更新密码
func (s *UserService) Password(id int64, user request.UserPassword) error {
	var (
		m  model.User
		db = s.DB(&m)
	)

	user.Password = pkg.BcryptHash(user.Password)

	return db.Model(&m).Where("id = ?", id).Updates(user).Error
}
