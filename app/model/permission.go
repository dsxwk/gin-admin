package model

const TableNamePermission = "permission"

// Permission 权限表
type Permission struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Key       string     `gorm:"column:key;not null;type:varchar(130);comment:权限标识" json:"key" form:"key"`
	Method    string     `gorm:"column:method;not null;type:varchar(20);comment:请求方式" json:"method" form:"method"`
	Uri       string     `gorm:"column:uri;not null;type:varchar(100);comment:路由地址" json:"uri" form:"uri"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*Permission) TableName() string {
	return TableNamePermission
}

// Connection 数据库连接名称
func (m *Permission) Connection() string {
	return "mysql"
}
