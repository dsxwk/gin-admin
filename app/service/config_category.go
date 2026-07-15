package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
)

type ConfigCategoryService struct {
	base.BaseService
}

// List 列表
func (s *ConfigCategoryService) List(req request.ConfigCategory) (pageData request.PageData, err error) {
	var (
		m      model.ConfigCategory
		models []model.ConfigCategory
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m)

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("id DESC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
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

// Create 创建
func (s *ConfigCategoryService) Create(req request.ConfigCategory) (request.ConfigCategory, error) {
	var (
		m  model.ConfigCategory
		db = s.DB(&m)
	)

	m = model.ConfigCategory{
		Name: req.Name,
	}

	err := db.Model(&m).Create(&m).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

// Update 更新
func (s *ConfigCategoryService) Update(id int64, data map[string]interface{}) (err error) {
	return s.Updates(&model.ConfigCategory{}, id, data)
}

// Detail 详情
func (s *ConfigCategoryService) Detail(id int64) (m model.ConfigCategory, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Delete 删除
func (s *ConfigCategoryService) Delete(id int64) (err error) {
	var (
		m  model.ConfigCategory
		db = s.DB(&m)
	)

	tx := db.Begin()

	err = tx.Model(&m).Delete(&m, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.SystemConfig{}).
		Where("config_category_id = ?", id).
		Delete(&model.SystemConfig{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
