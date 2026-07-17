package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
)

type ImportRecordsService struct {
	base.BaseService
}

// List 列表
func (s *ImportRecordsService) List(req request.ImportRecords) (pageData request.PageData, err error) {
	var (
		m      model.ImportRecords
		models []model.ImportRecords
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

// Delete 删除
func (s *ImportRecordsService) Delete(id int64) (err error) {
	var (
		m  model.ImportRecords
		db = s.DB(&m)
	)

	err = db.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}

	return nil
}
