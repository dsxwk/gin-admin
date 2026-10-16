package base

import (
	"context"
	"fmt"
	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/cache"
	"gin/pkg/serviceprovider/orm"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service 基础服务
type Service struct{}

// Model 数据库模型接口
type Model = orm.Model

// getDB 获取数据库连接
func (s *Service) getDB(ctx context.Context, model Model) *gorm.DB {
	conn := model.Connection()
	if conn != "" {
		return facade.DB(conn).WithContext(ctx)
	}
	return facade.DB().WithContext(ctx)
}

// DB 获取数据库连接
func (s *Service) DB(ctx context.Context, model Model) *gorm.DB {
	return s.getDB(ctx, model)
}

// Search 搜索扩展方法
func (s *Service) Search(db *gorm.DB, model any, conditions map[string]any) *gorm.DB {
	if len(conditions) == 0 {
		return db
	}

	whereSql, args, err := orm.BuildCondition(db, model, conditions)
	if err != nil {
		_ = db.AddError(fmt.Errorf("构建搜索条件失败:%w", err))
		return db
	}

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}
	return db
}

// WithoutSQLLog 临时关闭sql输出
func (s *Service) WithoutSQLLog(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{Logger: logger.Discard})
}

// Cache 获取缓存实例并绑定请求上下文
// 使用示例:
//
//	s.Cache(ctx).Set("key", value, 5*time.Minute)
//	s.Cache(ctx, "redis").Get("key")
func (s *Service) Cache(ctx context.Context, cacheType ...string) *cache.CacheProxy {
	return facade.Cache(cacheType...).WithContext(ctx)
}

// UpdateByID 根据ID更新数据
func (s *Service) UpdateByID(ctx context.Context, m Model, id int64, data map[string]any) error {
	db := s.DB(ctx, m)
	rows := model.FilterFields(db, m, data)
	if rows == nil {
		rows = make(map[string]any)
	}
	rows[model.UpdatedField] = time.Now()

	result := db.Model(m).Where("id = ?", id).Updates(rows)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var count int64
	if err := db.Model(m).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errcode.NotFound().WithMsg("数据不存在")
	}
	return nil
}
