package middleware

import (
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/errcode"

	"github.com/gin-gonic/gin"
)

type Permission struct {
	base.BaseMiddleware
}

// Handle 权限中间件
func (s Permission) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Request.Method + ":" + c.FullPath()

		if !s.hasPermission(c, code) {
			s.Response.Error(c, errcode.Forbidden().WithMsg("No permission"))
			return
		}
		c.Next()
	}
}

func (s Permission) hasPermission(c *gin.Context, code string) bool {
	userID := c.GetInt64(ctxkey.UserIdKey)
	ctx := c.Request.Context()
	redisKey := fmt.Sprintf("permission:user:%d", userID)

	redisCache := facade.Cache("redis").Redis().WithContext(ctx)

	// 检查Redis中权限集合是否存在
	exists, err := redisCache.Exists(redisKey)
	if err != nil {
		// Redis异常回退数据库查询
		return checkFromDB(userID, code)
	}

	if exists > 0 {
		// Redis有缓存直接查集合
		isMember, _err := redisCache.SIsMember(redisKey, code)
		if _err != nil {
			return checkFromDB(userID, code)
		}
		return isMember
	}

	// Redis无缓存从数据库加载
	permissions, err := loadUserPermissions(userID)
	if err != nil {
		return false
	}

	// 写入Redis缓存
	if len(permissions) > 0 {
		members := make([]any, len(permissions))
		for i, p := range permissions {
			members[i] = p
		}
		_ = redisCache.SAdd(redisKey, members...)
	}

	// 检查当前权限
	for _, p := range permissions {
		if p == code {
			return true
		}
	}
	return false
}

// checkFromDB 直接从数据库查询单个权限
func checkFromDB(userID int64, code string) bool {
	var count int64
	facade.DB("mysql").Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permission.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permission.key = ?", code).
		Where("permission.deleted_at IS NULL").
		Where("role_permissions.deleted_at IS NULL").
		Where("user_roles.deleted_at IS NULL").
		Count(&count)
	return count > 0
}

// loadUserPermissions 从数据库加载用户所有权限并缓存到Redis
func loadUserPermissions(userID int64) ([]string, error) {
	var permissions []model.Permission
	err := facade.DB("mysql").Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permission.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permission.deleted_at IS NULL").
		Where("role_permissions.deleted_at IS NULL").
		Where("user_roles.deleted_at IS NULL").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(permissions))
	for _, p := range permissions {
		keys = append(keys, p.Key)
	}
	return keys, nil
}
