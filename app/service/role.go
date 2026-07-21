package service

import (
	"errors"
	"fmt"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gorm.io/gorm"
	"time"
)

type RoleService struct {
	base.BaseService
}

// List 列表
func (s *RoleService) List(req request.Roles) (pageData request.PageData, err error) {
	var (
		m      model.Roles
		models []model.Roles
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).
		Model(&m).
		Preload("UserRoles").
		Preload("UserRoles.User").
		Preload("RoleMenus")

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
func (s *RoleService) Detail(id int64) (m model.Roles, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).
		Preload("UserRoles").
		Preload("UserRoles.User").
		Preload("RoleMenus").
		Preload("RolePermissions").
		First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Create 创建
func (s *RoleService) Create(req request.Roles) (m model.Roles, err error) {
	var (
		count           int64
		db              = s.DB(&m)
		roleMenus       []model.RoleMenus
		userRoles       []model.UserRoles
		rolePermissions []model.RolePermissions
	)

	// 校验角色名是否重复
	err = db.Model(&m).Where("name = ?", req.Name).Count(&count).Error
	if err != nil {
		return m, err
	}
	if count > 0 {
		return m, errors.New("角色名已存在")
	}

	m = model.Roles{
		Name:   req.Name,
		Desc:   req.Desc,
		Status: req.Status,
	}

	tx := db.Begin()

	err = tx.Model(&m).Create(&m).Error
	if err != nil {
		tx.Rollback()
		return m, err
	}

	if req.UserRoles != nil {
		for _, v := range req.UserRoles {
			userRoles = append(userRoles, model.UserRoles{
				UserID: v.UserId,
				RoleID: m.ID,
				Name:   v.Name,
			})
		}
		err = tx.Model(&model.UserRoles{}).Create(&userRoles).Error
		if err != nil {
			tx.Rollback()
			return m, err
		}
	}

	if req.RoleMenus != nil {
		for _, v := range req.RoleMenus {
			roleMenus = append(roleMenus, model.RoleMenus{
				RoleId: m.ID,
				MenuId: v.MenuId,
				Name:   m.Name,
			})
		}
		err = tx.Model(&model.RoleMenus{}).Create(&roleMenus).Error
		if err != nil {
			tx.Rollback()
			return m, err
		}
	}

	if req.RolePermissions != nil {
		for _, v := range req.RolePermissions {
			rolePermissions = append(rolePermissions, model.RolePermissions{
				RoleId:       m.ID,
				PermissionId: v.PermissionId,
			})
		}
		err = tx.Model(&model.RolePermissions{}).Create(&rolePermissions).Error
		if err != nil {
			tx.Rollback()
			return m, err
		}
	}

	tx.Commit()

	// 同步权限到Redis permission:user:{userId}
	if req.RolePermissions != nil {
		_ = s.syncRolePermissionsToRedis(m.ID)
	}

	return m, nil
}

// Update 更新
func (s *RoleService) Update(id int64, data map[string]interface{}) (err error) {
	var (
		count int64
		db    = s.DB(&model.Roles{})
	)

	// 校验角色名是否重复
	err = db.Model(&model.Roles{}).Where("name = ? AND id <> ?", data["name"], id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名已存在")
	}

	roleMenus, ok := data["roleMenus"].([]interface{})
	if !ok {
		roleMenus = []interface{}{}
	}

	userRoles, ok := data["userRoles"].([]interface{})
	if !ok {
		userRoles = []interface{}{}
	}

	rolePermissions, hasRolePermissions := data["rolePermissions"].([]interface{})
	if !hasRolePermissions {
		rolePermissions = []interface{}{}
	}

	tx := db.Begin()

	rows := model.FilterFields(s.DB(&model.Roles{}), model.Roles{}, data)
	rows[model.UpdatedField] = time.Now()

	err = tx.Model(&model.Roles{}).Where("id = ?", id).Updates(rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(userRoles) > 0 {
		// 删除旧关联
		err = tx.Model(&model.UserRoles{}).
			Where("role_id = ?", id).
			Delete(&model.UserRoles{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		var newUserRoles []model.UserRoles
		for _, item := range userRoles {
			userRoleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newUserRoles = append(newUserRoles, model.UserRoles{
				UserID: int64(userRoleMap["userId"].(float64)),
				RoleID: id,
				Name:   userRoleMap["name"].(string),
			})
		}

		if len(newUserRoles) > 0 {
			err = tx.Model(&model.UserRoles{}).Create(&newUserRoles).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 更新角色菜单关联
	if len(roleMenus) > 0 {
		// 删除旧关联
		err = tx.Model(&model.RoleMenus{}).
			Where("role_id = ?", id).
			Delete(&model.RoleMenus{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		var newRoleMenus []model.RoleMenus
		for _, item := range roleMenus {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newRoleMenus = append(newRoleMenus, model.RoleMenus{
				MenuId: int64(roleMap["menuId"].(float64)),
				RoleId: id,
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

	if hasRolePermissions {
		err = tx.Model(&model.RolePermissions{}).
			Where("role_id = ?", id).
			Delete(&model.RolePermissions{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		var newRolePermissions []model.RolePermissions
		for _, item := range rolePermissions {
			roleMap, _ok := item.(map[string]interface{})
			if !_ok {
				continue
			}

			newRolePermissions = append(newRolePermissions, model.RolePermissions{
				RoleId:       id,
				PermissionId: int64(roleMap["permissionId"].(float64)),
			})
		}

		if len(newRolePermissions) > 0 {
			err = tx.Model(&model.RolePermissions{}).Create(&newRolePermissions).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	// 同步权限到Redis permission:user:{userId}
	if hasRolePermissions {
		_ = s.syncRolePermissionsToRedis(id)
	}

	return nil
}

// Delete 删除
func (s *RoleService) Delete(id int64) (err error) {
	var (
		m  model.Roles
		db = s.DB(&m)
	)

	// 该角色的用户(事务前查询,事务后会删掉关联)
	var userRoles []model.UserRoles
	db.Model(&model.UserRoles{}).Where("role_id = ?", id).Find(&userRoles)

	userIDs := make([]int64, 0, len(userRoles))
	seen := make(map[int64]bool)
	for _, ur := range userRoles {
		if !seen[ur.UserID] {
			userIDs = append(userIDs, ur.UserID)
			seen[ur.UserID] = true
		}
	}

	tx := db.Begin()

	err = tx.Model(&m).Delete(&m, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.RoleMenus{}).
		Where("role_id = ?", id).
		Delete(&model.RoleMenus{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.UserRoles{}).
		Where("role_id = ?", id).
		Delete(&model.UserRoles{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.RolePermissions{}).
		Where("role_id = ?", id).
		Delete(&model.RolePermissions{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// 重建受影响用户的Redis权限集合
	for _, userID := range userIDs {
		_ = s.rebuildUserPermissions(db, userID)
	}

	return nil
}

// syncRolePermissionsToRedis 将角色的权限同步到其所有用户的Redis集合
func (s *RoleService) syncRolePermissionsToRedis(roleID int64) error {
	db := s.DB(&model.Roles{})

	// 查询该角色的所有用户
	var userRoles []model.UserRoles
	if err := db.Model(&model.UserRoles{}).Where("role_id = ?", roleID).Find(&userRoles).Error; err != nil {
		return err
	}
	if len(userRoles) == 0 {
		return nil
	}

	// 收集所有userId并去重
	userIDs := make([]int64, 0, len(userRoles))
	seen := make(map[int64]bool)
	for _, ur := range userRoles {
		if !seen[ur.UserID] {
			userIDs = append(userIDs, ur.UserID)
			seen[ur.UserID] = true
		}
	}

	// 为每个用户重建权限集合
	for _, userID := range userIDs {
		if err := s.rebuildUserPermissions(db, userID); err != nil {
			return err
		}
	}

	return nil
}

// userPermissionRow 用户权限查询结果
type userPermissionRow struct {
	UserID int64  `gorm:"column:user_id"`
	Key    string `gorm:"column:key"`
}

// SyncAllUserPermissions 全量同步所有用户权限到Redis(仅同步role_permissions表中已有的权限)
func (s *RoleService) SyncAllUserPermissions() error {
	db := s.DB(&model.Roles{})

	// 查所有有角色的用户
	var allUserIDs []int64
	if err := db.Model(&model.UserRoles{}).Distinct().Pluck("user_id", &allUserIDs).Error; err != nil {
		return err
	}
	if len(allUserIDs) == 0 {
		return nil
	}

	// 查所有用户的权限(仅角色权限表中已分配的)
	var rows []userPermissionRow
	db.Model(&model.Permission{}).
		Select("user_roles.user_id, permission.key").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permission.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("permission.deleted_at IS NULL").
		Where("role_permissions.deleted_at IS NULL").
		Where("user_roles.deleted_at IS NULL").
		Find(&rows)

	// 按用户分组：先初始化所有用户为空
	userPerms := make(map[int64][]string, len(allUserIDs))
	for _, uid := range allUserIDs {
		userPerms[uid] = nil
	}
	// 填充有权限的用户
	for _, row := range rows {
		userPerms[row.UserID] = append(userPerms[row.UserID], row.Key)
	}

	// Redis Pipeline 批量写入
	redisCache := s.Cache("redis").Redis()
	pipe := redisCache.Pipeline()
	for userID, keys := range userPerms {
		redisKey := fmt.Sprintf("permission:user:%d", userID)
		pipe.Del(s.Ctx, redisKey)
		if len(keys) > 0 {
			members := make([]interface{}, len(keys))
			for i, k := range keys {
				members[i] = k
			}
			pipe.SAdd(s.Ctx, redisKey, members...)
		}
	}

	_, err := pipe.Exec(s.Ctx)
	return err
}

// rebuildUserPermissions 重建单个用户的Redis权限集合
func (s *RoleService) rebuildUserPermissions(db *gorm.DB, userID int64) error {
	var (
		m           model.Permission
		permissions []model.Permission
	)

	err := db.Model(&m).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permission.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permission.deleted_at IS NULL").
		Where("role_permissions.deleted_at IS NULL").
		Where("user_roles.deleted_at IS NULL").
		Find(&permissions).Error
	if err != nil {
		return err
	}

	key := fmt.Sprintf("permission:user:%d", userID)
	redisCache := s.Cache("redis").Redis()

	_ = redisCache.Delete(key)

	if len(permissions) > 0 {
		members := make([]interface{}, len(permissions))
		for i, p := range permissions {
			members[i] = p.Key
		}
		_ = redisCache.SAdd(key, members...)
	}

	return nil
}
