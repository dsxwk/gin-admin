package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
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
