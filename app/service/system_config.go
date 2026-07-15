package service

import (
	"errors"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"time"
)

type SystemConfigService struct {
	base.BaseService
}

// List 列表
func (s *SystemConfigService) List(req request.SystemConfig) (pageData request.PageData, err error) {
	var (
		m      model.SystemConfig
		models []model.SystemConfig
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m).Preload("ConfigCategory")

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

// UpdateConfig 保存配置
func (s *SystemConfigService) UpdateConfig(data map[string]interface{}) (err error) {
	var (
		m  model.SystemConfig
		db = s.DB(&m)
	)

	rows, ok := data["list"]
	if !ok || rows == nil {
		return nil
	}

	// 先断言为[]interface{}
	listInterface, ok := rows.([]interface{})
	if !ok {
		return errors.New("list格式错误")
	}

	if len(listInterface) == 0 {
		return nil
	}
	// 转换为[]map[string]interface{}
	list := make([]map[string]interface{}, 0, len(listInterface))
	for _, item := range listInterface {
		mapItem, _ok := item.(map[string]interface{})
		if !_ok {
			return errors.New("list 中的元素格式错误")
		}
		list = append(list, mapItem)
	}

	sql, values := model.BatchUpdateSql(db, &m, list, "id", nil)
	err = db.Exec(sql, values...).Error
	if err != nil {
		return err
	}

	return nil
}

// Create 创建
func (s *SystemConfigService) Create(req request.SystemConfig) (request.SystemConfig, error) {
	var (
		m  model.SystemConfig
		db = s.DB(&m)
	)

	m = model.SystemConfig{
		Key:              req.Key,
		Name:             req.Name,
		DefaultValue:     req.DefaultValue,
		OptionValue:      req.OptionValue,
		Type:             req.Type,
		ConfigCategoryId: req.ConfigCategoryId,
	}

	err := db.Model(&m).Create(&m).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

// Update 更新
func (s *SystemConfigService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		m  model.SystemConfig
		db = s.DB(&m)
	)

	rows := model.FilterFields(db, m, data)
	rows["updated_at"] = time.Now()

	err = db.Model(&m).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		return err
	}

	return nil
}

// Detail 详情
func (s *SystemConfigService) Detail(id int64) (m model.SystemConfig, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).Preload("ConfigCategory").First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Delete 删除
func (s *SystemConfigService) Delete(id int64) (err error) {
	var (
		m  model.SystemConfig
		db = s.DB(&m)
	)

	err = db.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}

	return nil
}
