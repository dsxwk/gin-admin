package base

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/cache"
	"gin/pkg/serviceprovider/orm"
	"time"

	"gorm.io/gorm"
)

type BaseService struct {
}

type Model interface {
	TableName() string
}

// HasConnection 检查是否有连接方法的接口
type HasConnection interface {
	Connection() string
}

// getDB 获取数据库连接(带连接名判断)
func (s *BaseService) getDB(ctx context.Context, model Model) *gorm.DB {
	if connModel, ok := model.(HasConnection); ok {
		conn := connModel.Connection()
		if conn != "" {
			return facade.DB(conn).WithContext(ctx)
		}
	}
	return facade.DB().WithContext(ctx)
}

// DB 获取数据库连接
func (s *BaseService) DB(ctx context.Context, model Model) *gorm.DB {
	return s.getDB(ctx, model)
}

// Search 搜索扩展方法
func (s *BaseService) Search(db *gorm.DB, model any, conditions map[string]any) *gorm.DB {
	if conditions == nil || len(conditions) == 0 {
		return db
	}

	whereSql, args, err := orm.BuildCondition(db, model, conditions)
	if err != nil {
		return db
	}

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}
	return db
}

// Cache 获取缓存实例并绑定请求上下文
// 使用示例:
//
//	s.Cache(ctx).Set("key", value, 5*time.Minute)
//	s.Cache(ctx, "redis").Get("key")
func (s *BaseService) Cache(ctx context.Context, cacheType ...string) *cache.CacheProxy {
	return facade.Cache(cacheType...).WithContext(ctx)
}

// Updates 公共更新方法
func (s *BaseService) Updates(ctx context.Context, m Model, id int64, data map[string]any) error {
	db := s.DB(ctx, m)
	rows := model.FilterFields(db, m, data)
	rows[model.UpdatedField] = time.Now()
	return db.Model(m).Where("id = ?", id).Updates(rows).Error
}
