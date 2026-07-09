package model

const TableNameConfigCategory = "config_category"

// ConfigCategory 配置分类表
type ConfigCategory struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Name      string     `gorm:"column:name;not null;type:varchar(50);comment:分类名称" json:"name" form:"name"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*ConfigCategory) TableName() string {
	return TableNameConfigCategory
}

// Connection 数据库连接名称
func (m *ConfigCategory) Connection() string {
	return "mysql"
}
