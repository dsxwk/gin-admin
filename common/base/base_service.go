package base

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/orm"
	"gorm.io/gorm"
)

type BaseService struct {
	Context
}

type Model interface {
	TableName() string
}

// HasConnection 检查是否有连接方法的接口
type HasConnection interface {
	Connection() string
}

func (s *BaseService) WithContext(ctx context.Context) *BaseService {
	s.Set(ctx)
	return s
}

// getDB 获取数据库连接(带连接名判断)
func (s *BaseService) getDB(model Model) *gorm.DB {
	if connModel, ok := model.(HasConnection); ok {
		conn := connModel.Connection()
		if conn != "" {
			return facade.DB(conn).WithContext(s.Ctx)
		}
	}
	return facade.DB().WithContext(s.Ctx)
}

// DB 获取数据库连接,连接无效时自动重连
func (s *BaseService) DB(model Model) *gorm.DB {
	db := s.getDB(model)

	// 检查连接是否有效
	sqlDB, err := db.DB()
	if err == nil {
		if err = sqlDB.Ping(); err != nil {
			facade.Log().Warn("数据库连接无效,重新连接")
			db = s.getDB(model)
		}
	}

	return db
}

// Search 搜索扩展方法
func (s *BaseService) Search(db *gorm.DB, model any, conditions map[string]interface{}) *gorm.DB {
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
