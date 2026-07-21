package service

import (
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"strings"
)

type PermissionService struct {
	base.BaseService
}

// List 权限列表
func (s *PermissionService) List(req request.Permission) (pageData request.PageData, err error) {
	var (
		m      model.Permission
		models []model.Permission
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).Model(&m)

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

// SyncRoutePermissions 同步路由权限
// permissionKeys METHOD:PATH (如 GET:/api/v1/user)
func (s *PermissionService) SyncRoutePermissions(permissionKeys []string) error {
	if len(permissionKeys) == 0 {
		return nil
	}

	var (
		m      model.Permission
		models []model.Permission
		db     = s.DB(&m)
	)

	// 解析路由Key
	routeMap := make(map[string][2]string, len(permissionKeys))
	for _, key := range permissionKeys {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			continue
		}
		routeMap[key] = [2]string{parts[0], parts[1]}
	}

	routeKeyList := make([]string, 0, len(routeMap))
	for k := range routeMap {
		routeKeyList = append(routeKeyList, k)
	}

	// 查询数据库中已有的权限
	db.Model(&m).Where("`key` IN ?", routeKeyList).Find(&models)
	existingMap := make(map[string]bool, len(models))
	for _, p := range models {
		existingMap[p.Key] = true
	}

	// 只插入不存在的
	var data []model.Permission
	for k, v := range routeMap {
		if !existingMap[k] {
			data = append(data, model.Permission{
				Key:    k,
				Method: v[0],
				Uri:    v[1],
			})
		}
	}
	if len(data) > 0 {
		db.Model(&m).Create(&data)
	}

	// 清理已删除的路由权限
	db.Model(&model.Permission{}).
		Where("`key` NOT IN ?", routeKeyList).
		Delete(&model.Permission{})

	return nil
}
