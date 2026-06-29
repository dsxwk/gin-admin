package service

import (
	"errors"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"strings"
	"time"
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
	err = s.DB(&m).
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
		m         model.Menu
		meta      model.MenuMeta
		roleMenus []model.RoleMenus
		db        = s.DB(&model.Menu{})
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

	db = db.Begin()
	err := db.Model(&model.Menu{}).Create(&m).Error
	if err != nil {
		db.Rollback()
		return req, err
	}

	meta = model.MenuMeta{
		MenuId:      m.ID,
		Title:       req.Meta.Title,
		Icon:        req.Meta.Icon,
		IsHide:      req.Meta.IsHide,
		IsKeepAlive: req.Meta.IsKeepAlive,
		IsAffix:     req.Meta.IsAffix,
		IsLink:      req.Meta.IsLink,
		IsIframe:    req.Meta.IsIframe,
	}

	err = db.Model(&model.MenuMeta{}).Create(&meta).Error
	if err != nil {
		db.Rollback()
		return req, err
	}

	if req.RoleMenus != nil {
		for _, v := range req.RoleMenus {
			roleMenus = append(roleMenus, model.RoleMenus{
				RoleId: v.RoleId,
				MenuId: m.ID,
			})
		}
		err = db.Model(&model.RoleMenus{}).Create(&roleMenus).Error
		if err != nil {
			db.Rollback()
			return req, err
		}
	}

	db.Commit()

	return req, nil
}

// Update 更新
func (s *MenuService) Update(id int64, data map[string]interface{}) error {
	var (
		err error
		db  = s.DB(&model.Menu{})
	)

	meta, ok := data["meta"].(map[string]interface{})
	if !ok {
		return errors.New("meta数据格式错误")
	}

	roleMenusData, ok := data["roleMenus"].([]interface{})
	if !ok {
		roleMenusData = []interface{}{}
	}

	tx := db.Begin()

	rows := model.FilterFields(db, model.Menu{}, data)
	rows["updated_at"] = time.DateTime

	err = tx.Model(&model.Menu{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	metaRows := model.FilterFields(db, model.MenuMeta{}, meta)
	metaRows["updated_at"] = time.DateTime
	err = tx.Model(&model.MenuMeta{}).
		Where("id = ?", meta["id"]).
		Updates(metaRows).Error

	// 更新角色菜单关联
	if len(roleMenusData) > 0 {
		// 删除旧关联
		err = tx.Model(&model.RoleMenus{}).
			Where("menu_id = ?", id).
			Delete(&model.RoleMenus{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		var newRoleMenus []model.RoleMenus
		for _, item := range roleMenusData {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newRoleMenus = append(newRoleMenus, model.RoleMenus{
				MenuId: id,
				RoleId: int64(roleMap["roleId"].(float64)),
			})
		}

		if len(newRoleMenus) > 0 {
			err = tx.Model(&model.RoleMenus{}).Create(&newRoleMenus).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// Delete 删除
func (s *MenuService) Delete(menuId int64) error {
	var (
		err           error
		menuActionIds []int64
		db            = s.DB(&model.Menu{})
		roleMenu      model.RoleMenus
		menuMeta      model.MenuMeta
		menuAction    model.MenuActions
		roleAction    model.RoleActions
	)

	tx := db.Begin()

	err = tx.Delete(&model.Menu{}, menuId).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.RoleMenus{}).
		Where("menu_id = ?", menuId).
		Delete(&roleMenu).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.MenuMeta{}).
		Where("menu_id = ?", menuId).
		Delete(&menuMeta).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.MenuActions{}).
		Where("menu_id = ?", menuId).
		Pluck("id", &menuActionIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.MenuActions{}).
		Where("menu_id = ?", menuId).
		Delete(&model.MenuActions{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(menuActionIds) > 0 {
		err = tx.Model(&model.MenuActions{}).
			Where("id IN (?)", menuActionIds).
			Delete(&menuAction).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Model(&model.RoleActions{}).
			Where("action_id IN (?)", menuActionIds).
			Delete(&roleAction).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

// Action 菜单功能
func (s *MenuService) Action(req request.MenuActions) (pageData request.PageData, err error) {
	var (
		m      []model.MenuActions
		action model.MenuActions
		db     = s.DB(&action)
	)

	// 搜索
	db = s.Search(db, req.Search)

	err = db.Preload("Parent").
		Preload("RoleActions").
		Where("menu_id = ?", req.Id).
		Order("sort asc").
		Find(&m).
		Error

	pageData.List = action.GetTree(m)

	return pageData, nil
}

// ActionDetail 菜单功能详情
func (s *MenuService) ActionDetail(id int64) (m model.MenuActions, err error) {
	err = s.DB(&model.MenuActions{}).
		Preload("Parent").
		Preload("RoleActions").
		First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// CreateAction 菜单功能创建
func (s *MenuService) CreateAction(req request.MenuActions) (err error) {
	var (
		menuAction  model.MenuActions
		roleActions []model.RoleActions
		db          = s.DB(&model.MenuActions{})
	)

	menuAction = model.MenuActions{
		Pid:       req.Pid,
		MenuId:    req.MenuId,
		Type:      req.Type,
		BtnType:   req.BtnType,
		BtnStyle:  req.BtnStyle,
		BtnSize:   req.BtnSize,
		IsConfirm: req.IsConfirm,
		Label:     req.Label,
		AuthValue: req.AuthValue,
		IsLink:    req.IsLink,
		Sort:      req.Sort,
	}

	tx := db.Begin()

	err = tx.Create(&menuAction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range req.RoleActions {
		roleActions = append(
			roleActions,
			model.RoleActions{
				RoleId:   v.RoleId,
				Name:     v.Name,
				ActionId: menuAction.Id,
			},
		)
	}

	err = tx.Model(&model.RoleActions{}).Create(&roleActions).Error

	tx.Commit()

	return nil
}

// UpdateAction 菜单功能更新
func (s *MenuService) UpdateAction(actionId int64, data map[string]interface{}) (err error) {
	var (
		db = s.DB(&model.MenuActions{})
	)

	rows := model.FilterFields(db, model.MenuActions{}, data)
	roleActionData, ok := data["roleActions"].([]interface{})
	if !ok {
		roleActionData = []interface{}{}
	}

	tx := db.Begin()

	err = tx.Where("id = ?", actionId).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除旧关联
	err = tx.Model(&model.RoleActions{}).
		Where("action_id = ?", actionId).
		Delete(&model.RoleMenus{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新角色功能关联
	if len(roleActionData) > 0 {
		// 创建新关联
		var newRoleActions []model.RoleActions
		for _, item := range roleActionData {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			name, _ := roleMap["name"].(string)

			newRoleActions = append(newRoleActions, model.RoleActions{
				RoleId:   int64(roleMap["roleId"].(float64)),
				Name:     name,
				ActionId: actionId,
			})
		}

		if len(newRoleActions) > 0 {
			err = tx.Model(&model.RoleActions{}).Create(&newRoleActions).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// DeleteAction 菜单功能删除
func (s *MenuService) DeleteAction(id int64) (err error) {
	var (
		roleAction model.RoleActions
		db         = s.DB(&model.MenuActions{})
	)

	tx := db.Begin()

	err = tx.Delete(&model.MenuActions{}, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.RoleActions{}).Where("action_id = ?", id).Delete(&roleAction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
