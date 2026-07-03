package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"time"
)

type ArticleService struct {
	base.BaseService
}

// List 列表
func (s *ArticleService) List(req request.Article) (pageData request.PageData, err error) {
	var (
		models []model.Article
		m      model.Article
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, req.Search)

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
func (s *ArticleService) Create(req request.Article) (request.Article, error) {
	var (
		m model.Article
	)

	m = model.Article{
		Uid:        req.Uid,
		Title:      req.Title,
		Content:    req.Content,
		CategoryId: req.CategoryId,
		DataSource: req.DataSource,
		IsPublish:  req.IsPublish,
		Tag:        &model.JsonValue{Data: req.Tag},
	}

	err := s.DB(&model.Article{}).Create(&m).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

// Update 更新
func (s *ArticleService) Update(id int64, data map[string]interface{}) (err error) {
	if pkg.HasKey(data, "tag") {
		data["tag"] = &model.JsonValue{Data: data["tag"]}
	}
	rows := model.FilterFields(s.DB(&model.Article{}), model.Article{}, data)
	rows["updated_at"] = time.DateTime

	err = s.DB(&model.Article{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		return err
	}

	return nil
}

// Detail 详情
func (s *ArticleService) Detail(id int64) (m model.Article, err error) {
	err = s.DB(&model.Article{}).First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Delete 删除
func (s *ArticleService) Delete(id int64) (err error) {
	var (
		m  model.Article
		db = s.DB(&m)
	)

	err = db.Delete(&m, id).Error
	if err != nil {
		return err
	}

	return nil
}
