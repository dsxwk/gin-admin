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
		m      model.Menu
		models []model.Menu
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).
		Model(&m).
		Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("MenuAction").
		Preload("RoleMenus")

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

// RoleMenu 角色菜单
func (s *MenuService) RoleMenu(req request.Menu) (tree []pkg.TreeNode, err error) {
	var (
		m      model.Menu
		models []model.Menu
		db     = s.DB(&m)
	)

	roleIds := strings.Split(req.RoleId, ",")

	err = db.Model(&m).
		Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("RoleMenus", "role_id IN ?", roleIds).
		Preload("MenuAction").
		Joins("LEFT JOIN role_menus ON menu.id = role_menus.menu_id").
		Where("role_menus.role_id IN ?", roleIds).
		Order("sort ASC").
		Group("menu.id").
		Find(&models).Error
	if err != nil {
		return tree, err
	}

	return m.GetTree(models), nil
}

// Detail 详情
func (s *MenuService) Detail(menuId int64) (m model.Menu, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).
		Preload("Meta").
		Preload("Meta.AuthBtnList").
		Preload("MenuAction").
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
		m          model.Menu
		meta       model.MenuMeta
		menuAction model.MenuActions
		roleMenus  []model.RoleMenus
		db         = s.DB(&m)
	)

	m = model.Menu{
		Type:   req.Type,
		Pid:    req.Pid,
		Name:   req.Name,
		Status: req.Status,
		Sort:   req.Sort,
	}

	db = db.Begin()
	err := db.Model(&m).Create(&m).Error
	if err != nil {
		db.Rollback()
		return req, err
	}

	if req.Meta.Title != "" {
		meta = model.MenuMeta{
			MenuId:      m.ID,
			Title:       req.Meta.Title,
			Icon:        req.Meta.Icon,
			Path:        req.Meta.Path,
			Redirect:    req.Meta.Redirect,
			Component:   req.Meta.Component,
			IsHide:      req.Meta.IsHide,
			IsKeepAlive: req.Meta.IsKeepAlive,
			IsAffix:     req.Meta.IsAffix,
			IsLink:      req.Meta.IsLink,
			IsIframe:    req.Meta.IsIframe,
		}

		err = db.Model(&meta).Create(&meta).Error
		if err != nil {
			db.Rollback()
			return req, err
		}
	}

	if req.MenuAction.Type != 0 {
		menuAction = model.MenuActions{
			MenuId:    m.ID,
			Type:      req.MenuAction.Type,
			AuthValue: req.MenuAction.AuthValue,
			BtnSize:   req.MenuAction.BtnSize,
			BtnStyle:  req.MenuAction.BtnStyle,
			BtnType:   req.MenuAction.BtnType,
			IsConfirm: req.MenuAction.IsConfirm,
			IsLink:    req.MenuAction.IsLink,
			Label:     req.MenuAction.Label,
		}

		err = db.Model(&menuAction).Create(&menuAction).Error
		if err != nil {
			db.Rollback()
			return req, err
		}
	}

	if req.RoleMenus != nil {
		for _, v := range req.RoleMenus {
			roleMenus = append(roleMenus, model.RoleMenus{
				RoleId: v.RoleId,
				MenuId: m.ID,
				Name:   v.Name,
			})
		}
		err = db.Model(&roleMenus).Create(&roleMenus).Error
		if err != nil {
			db.Rollback()
			return req, err
		}
	}

	db.Commit()

	return req, nil
}

// Update 更新
func (s *MenuService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		m          model.Menu
		db         = s.DB(&m)
		menuAction map[string]interface{}
		meta       map[string]interface{}
		ok         bool
	)

	if pkg.HasKey(data, "meta") {
		meta, ok = data["meta"].(map[string]interface{})
		if !ok {
			return errors.New("meta数据格式错误")
		}
	}

	if pkg.HasKey(data, "menuAction") {
		menuAction, ok = data["menuAction"].(map[string]interface{})
		if !ok {
			return errors.New("菜单功能数据格式错误")
		}
	}

	roleMenusData, ok := data["roleMenus"].([]interface{})
	if !ok {
		roleMenusData = []interface{}{}
	}

	tx := db.Begin()

	rows := model.FilterFields(db, m, data)
	rows["updated_at"] = time.DateTime

	err = tx.Model(&m).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(meta) > 0 {
		metaRows := model.FilterFields(db, model.MenuMeta{}, meta)
		metaRows["updated_at"] = time.DateTime
		err = tx.Model(&model.MenuMeta{}).
			Where("id = ?", meta["id"]).
			Updates(metaRows).Error
	}

	if len(menuAction) > 0 {
		menuActionRows := model.FilterFields(db, model.MenuActions{}, menuAction)
		menuActionRows["updated_at"] = time.DateTime
		err = tx.Model(&model.MenuActions{}).
			Where("menu_id = ?", menuAction["menuId"]).
			Updates(menuActionRows).Error
	}

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
				Name:   roleMap["name"].(string),
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
func (s *MenuService) Delete(menuId int64) (err error) {
	var (
		menuActionIds []int64
		db            = s.DB(&model.Menu{})
		roleMenu      model.RoleMenus
		menuMeta      model.MenuMeta
		menuAction    model.MenuActions
	)

	tx := db.Begin()

	err = tx.Model(&model.Menu{}).Delete(&model.Menu{}, menuId).Error
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
	}

	tx.Commit()

	return nil
}
