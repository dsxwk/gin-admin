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
		m    []model.Roles
		role model.Roles
		db   = s.DB(&role)
	)

	pageData.Page = req.Page
	pageData.PageSize = req.PageSize
	offset, limit := request.Pagination(req.Page, req.PageSize)
	// 搜索
	db = s.Search(db, req.Search)

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
		err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m
	}

	return pageData, nil
}

// Detail 详情
func (s *RoleService) Detail(id int64) (role model.Roles, err error) {
	err = s.DB(&model.Roles{}).First(&role, id).Error
	if err != nil {
		return role, err
	}

	return role, nil
}

// Create 创建
func (s *RoleService) Create(req request.Roles) (m model.Roles, err error) {
	var (
		count int64
		db    = s.DB(&model.Roles{})
	)

	// 校验角色名是否重复
	err = db.Where("name = ?", req.Name).Count(&count).Error
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

	err = db.Create(&m).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Update 更新
func (s *RoleService) Update(id int64, data map[string]interface{}) error {
	var (
		count int64
	)

	// 校验角色名是否重复
	err := s.DB(&model.Roles{}).Where("name = ? AND id <> ?", data["name"], id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名已存在")
	}

	rows := model.FilterFields(s.DB(&model.Roles{}), model.Roles{}, data)
	rows["updated_at"] = time.DateTime

	err = s.DB(&model.Roles{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (s *RoleService) Delete(id int64) error {
	var (
		err  error
		role model.Roles
		db   = s.DB(&role)
	)

	err = db.Delete(&role, id).Error
	if err != nil {
		return err
	}

	return nil
}
