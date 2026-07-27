package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
)

type OperatorLogService struct {
	base.BaseService
}

// List 列表
func (s *OperatorLogService) List(req request.OperatorLog) (pageData request.PageData, err error) {
	var (
		m      model.OperatorLog
		models []model.OperatorLog
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m).Preload("User")

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

// Detail 详情
func (s *OperatorLogService) Detail(id int64) (m model.OperatorLog, err error) {
	db := s.DB(&m)
	err = db.Model(&m).Preload("User").First(&m, id).Error
	if err != nil {
		return m, err
	}
	return m, nil
}

// Delete 删除
func (s *OperatorLogService) Delete(id int64) (err error) {
	var (
		m  model.OperatorLog
		db = s.DB(&m)
	)
	err = db.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchDelete 批量删除
func (s *OperatorLogService) BatchDelete(ids []int64) (err error) {
	var (
		m  model.OperatorLog
		db = s.DB(&m)
	)
	err = db.Model(&m).Delete(&m, ids).Error
	if err != nil {
		return err
	}
	return nil
}
