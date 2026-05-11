package base

import (
	"context"
	"gin/app/facade"
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

	return db.Model(model)
}
