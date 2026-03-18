package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg/container"
	"gin/pkg/gorm/search"
)

type MenuService struct {
	base.BaseService
}

// List 列表
func (s *MenuService) List(req request.Menu, _search map[string]interface{}) (pageData request.PageData, err error) {
	var (
		m          []model.Menu
		menu       model.Menu
		containers = container.Get(s.Ctx)
	)

	offset, limit := request.Pagination(req.Page, req.PageSize)

	db := containers.DB.Model(&menu)

	if _search != nil {
		whereSql, args, _err := search.BuildCondition(_search, db, menu)
		if _err != nil {
			return pageData, err
		}

		if whereSql != "" {
			db = db.Where(whereSql, args...)
		}
	}

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("id DESC").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = menu.GetTree(m)
	} else {
		err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m
	}

	return pageData, nil
}
