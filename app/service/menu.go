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
	db = db.Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("MenuActions").
		Preload("RoleMenus")

	if req.NotPage {
		err = db.Order("sort Asc").Find(&m).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = menu.GetTree(m)
	} else {
		err = db.Offset(offset).Limit(limit).Order("sort Asc").Find(&m).Error
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
		m    []model.Menu
		menu model.Menu
		db   = s.DB(&menu)
	)

	roleIds := strings.Split(req.RoleId, ",")

	err = db.
		Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("RoleMenus").
		Preload("MenuActions").
		Preload("MenuActions.RoleActions", "role_id IN ?", roleIds).
		Preload("MenuActions.Parent").
		Joins("LEFT JOIN role_menus ON menu.id = role_menus.menu_id").
		Where("role_menus.role_id IN ?", roleIds).
		Order("sort ASC").
		Group("menu.id").
		Find(&m).Error
	if err != nil {
		return tree, err
	}

	return menu.GetTree(m), nil
}

// Detail 详情
func (s *MenuService) Detail(menuId int64) (m model.Menu, err error) {
	err = s.DB(&model.Menu{}).
		Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("MenuActions").
		Preload("RoleMenus").
		Order("sort Asc").
		First(&m, menuId).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Create 新增菜单
func (s *MenuService) Create(req request.Menu) (request.Menu, error) {
	var (
		m    model.Menu
		meta model.MenuMeta
		db   = s.DB(&model.Menu{})
	)

	m = model.Menu{
		Pid:       req.Pid,
		Name:      req.Name,
		Path:      req.Path,
		Redirect:  req.Redirect,
		Component: req.Component,
		IsLink:    req.IsLink,
		Status:    req.Status,
		Sort:      req.Sort,
	}

	meta = model.MenuMeta{
		MenuId:      req.Mata.MenuId,
		Title:       req.Mata.Title,
		Icon:        req.Mata.Icon,
		IsHide:      req.Mata.IsHide,
		IsKeepAlive: req.Mata.IsKeepAlive,
		IsAffix:     req.Mata.IsAffix,
		IsLink:      req.Mata.IsLink,
		IsIframe:    req.Mata.IsIframe,
	}

	db = db.Begin()
	err := db.Model(&model.Menu{}).Create(&m).Error
	if err != nil {
		db.Rollback()
		return req, err
	}

	err = db.Model(&model.MenuMeta{}).Create(&meta).Error
	if err != nil {
		db.Rollback()
		return req, err
	}

	db.Commit()

	return req, nil
}

// Action 菜单功能
func (s *MenuService) Action(menuId int64) (pageData request.PageData, err error) {
	var (
		m      []model.MenuActions
		action model.MenuActions
		db     = s.DB(&action)
	)

	err = db.Preload("Parent").
		Preload("RoleActions").
		Where("menu_id = ?", menuId).
		Order("sort asc").
		Find(&m).
		Error

	pageData.List = action.GetTree(m)

	return pageData, nil
}
