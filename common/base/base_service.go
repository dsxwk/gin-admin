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

func (s *BaseService) DB(model Model) *gorm.DB {
	var db *gorm.DB
	if connModel, ok := model.(HasConnection); ok {
		conn := connModel.Connection()
		if conn != "" {
			db = facade.DB(conn).WithContext(s.Ctx)
		} else {
			db = facade.DB().WithContext(s.Ctx)
		}
	} else {
		// 默认连接
		db = facade.DB().WithContext(s.Ctx)
	}

	// 检查连接是否有效
	sqlDB, err := db.DB()
	if err == nil {
		if err = sqlDB.Ping(); err != nil {
			// 连接无效,重新获取
			facade.Log().Warn("数据库连接无效,重新连接")
			if connModel, ok := model.(HasConnection); ok {
				conn := connModel.Connection()
				if conn != "" {
					db = facade.DB(conn).WithContext(s.Ctx)
				} else {
					db = facade.DB().WithContext(s.Ctx)
				}
			} else {
				db = facade.DB().WithContext(s.Ctx)
			}
		}
	}

	return db.Model(model)
}

// Search 搜索扩展方法
func (s *BaseService) Search(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	if conditions == nil || len(conditions) == 0 {
		return db
	}

	whereSql, args, err := orm.BuildCondition(db, conditions)
	if err != nil {
		return db
	}

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}
	return db
}
