package service

import (
	"errors"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"time"
)

type RoleService struct {
	base.BaseService
}

// List 列表
func (s *RoleService) List(req request.Roles) (pageData request.PageData, err error) {
	var (
		m      model.Roles
		models []model.Roles
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m).Preload("RoleMenus")

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("id DESC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m
	} else {
		pageData.Page = req.Page
		pageData.PageSize = req.PageSize
		offset, limit := request.Pagination(req.Page, req.PageSize)

		err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	}

	return pageData, nil
}

// Detail 详情
func (s *RoleService) Detail(id int64) (m model.Roles, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).
		Preload("RoleMenus").
		First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Create 创建
func (s *RoleService) Create(req request.Roles) (m model.Roles, err error) {
	var (
		count     int64
		db        = s.DB(&m)
		roleMenus []model.RoleMenus
	)

	// 校验角色名是否重复
	err = db.Model(&m).Where("name = ?", req.Name).Count(&count).Error
	if err != nil {
		return m, err
	}
	if count > 0 {
		return m, errors.New("角色名已存在")
	}

	m = model.Roles{
		Name:   req.Name,
		Desc:   req.Desc,
		Status: req.Status,
	}

	tx := db.Begin()

	err = tx.Model(&m).Create(&m).Error
	if err != nil {
		tx.Rollback()
		return m, err
	}

	if req.RoleMenus != nil {
		for _, v := range req.RoleMenus {
			roleMenus = append(roleMenus, model.RoleMenus{
				RoleId: m.ID,
				MenuId: v.MenuId,
				Name:   m.Name,
			})
		}
		err = tx.Model(&model.RoleMenus{}).Create(&roleMenus).Error
		if err != nil {
			tx.Rollback()
			return m, err
		}
	}

	tx.Commit()

	return m, nil
}

// Update 更新
func (s *RoleService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		count int64
		db    = s.DB(&model.Roles{})
	)

	// 校验角色名是否重复
	err = db.Model(&model.Roles{}).Where("name = ? AND id <> ?", data["name"], id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名已存在")
	}

	roleMenusData, ok := data["roleMenus"].([]interface{})
	if !ok {
		roleMenusData = []interface{}{}
	}

	tx := db.Begin()

	rows := model.FilterFields(s.DB(&model.Roles{}), model.Roles{}, data)
	rows["updated_at"] = time.DateTime

	err = tx.Model(&model.Roles{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新角色菜单关联
	if len(roleMenusData) > 0 {
		// 删除旧关联
		err = tx.Model(&model.RoleMenus{}).
			Where("role_id = ?", id).
			Delete(&model.RoleMenus{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		var newRoleMenus []model.RoleMenus
		for _, item := range roleMenusData {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newRoleMenus = append(newRoleMenus, model.RoleMenus{
				MenuId: int64(roleMap["menuId"].(float64)),
				RoleId: id,
			})
		}

		if len(newRoleMenus) > 0 {
			err = tx.Model(&model.RoleMenus{}).Create(&newRoleMenus).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// Delete 删除
func (s *RoleService) Delete(id int64) (err error) {
	var (
		m  model.Roles
		db = s.DB(&m)
	)

	tx := db.Begin()

	err = tx.Model(&m).Delete(&m, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.RoleMenus{}).
		Where("role_id = ?", id).
		Delete(&model.RoleMenus{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
