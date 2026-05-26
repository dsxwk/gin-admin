package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"strings"
)

type MenuService struct {
	base.BaseService
}

// List 列表
func (s *MenuService) List(req request.Menu) (pageData request.PageData, err error) {
	var (
		m    []model.Menu
		menu model.Menu
		db   = s.DB(&menu)
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

// RoleMenu 角色菜单
func (s *MenuService) RoleMenu(req request.Menu) (tree []pkg.TreeNode, err error) {
	var (
		m     []model.Menu
		menu  model.Menu
		db    = s.DB(&menu)
		total int64
	)

	roleIds := strings.Split(req.RoleId, ",")

	// 搜索
	db = s.Search(db, req.Search)

	err = db.Count(&total).Error
	if err != nil {
		return tree, err
	}

	err = db.
		Preload("RoleMenus").
		Preload("MenuActions").
		Preload("MenuActions.RoleActions", "role_id IN ?", roleIds).
		Joins("LEFT JOIN role_menus ON menu.id = role_menus.menu_id").
		Where("role_menus.role_id IN ?", roleIds).
		Order("sort asc").
		Group("menu.id").
		Find(&m).Error
	if err != nil {
		return tree, err
	}

	return menu.GetTree(m), nil
}
