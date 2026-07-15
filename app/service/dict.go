package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"time"
)

type DictService struct {
	base.BaseService
}

// List 列表
func (s *DictService) List(req request.Dict) (pageData request.PageData, err error) {
	var (
		m      model.Dict
		models []model.Dict
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m)

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("sort Asc").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m.GetTree(models)
	} else {
		pageData.Page = req.Page
		pageData.PageSize = req.PageSize
		offset, limit := request.Pagination(req.Page, req.PageSize)

		err = db.Offset(offset).Limit(limit).Order("sort Asc").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	}

	return pageData, nil
}

// Create 创建
func (s *DictService) Create(req request.Dict) (request.Dict, error) {
	var (
		m  model.Dict
		db = s.DB(&m)
	)

	m = model.Dict{
		Pid:    req.Pid,
		Name:   req.Name,
		Title:  req.Title,
		Value:  req.Value,
		Status: req.Status,
		Sort:   req.Sort,
		Extend: &model.JsonValue{Data: req.Extend},
		Desc:   req.Desc,
	}

	err := db.Model(&m).Create(&m).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

// Update 更新
func (s *DictService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		m  model.Dict
		db = s.DB(&m)
	)

	if pkg.HasKey(data, "extend") {
		data["extend"] = &model.JsonValue{Data: data["extend"]}
	}
	rows := model.FilterFields(db, m, data)
	rows["updated_at"] = time.Now()

	err = db.Model(&m).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		return err
	}

	return nil
}

// Detail 详情
func (s *DictService) Detail(id int64) (m model.Dict, err error) {
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
func (s *DictService) Delete(id int64) (err error) {
	var (
		m  model.Dict
		db = s.DB(&m)
	)

	err = db.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}

	return nil
}
