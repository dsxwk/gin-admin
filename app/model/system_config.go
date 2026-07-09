package model

const TableNameSystemConfig = "system_config"

// SystemConfig 系统配置表
type SystemConfig struct {
	ID               int64           `gorm:"column:id;primaryKey;autoIncrement;not null;type:smallint(5) unsigned;comment:id" json:"id" form:"id"`
	Key              string          `gorm:"column:key;not null;type:varchar(60);comment:标识" json:"key" form:"key"`
	Name             string          `gorm:"column:name;not null;type:varchar(60);comment:名称" json:"name" form:"name"`
	DefaultValue     string          `gorm:"column:default_value;not null;type:varchar(200);comment:默认值" json:"defaultValue" form:"defaultValue"`
	OptionValue      string          `gorm:"column:option_value;not null;type:varchar(200);comment:可选值" json:"optionValue" form:"optionValue"`
	Type             int64           `gorm:"column:type;not null;default:1;type:tinyint(1) unsigned;comment:配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件" json:"type" form:"type"`
	ConfigCategoryId int64           `gorm:"column:config_category_id;not null;default:0;type:tinyint(1) unsigned;comment:配置分类Id" json:"configCategoryId" form:"configCategoryId"`
	ConfigCategory   *ConfigCategory `gorm:"foreignKey:config_category_id;references:id" json:"configCategory"`
	CreatedAt        *DateTime       `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt        *DateTime       `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt        *DeletedAt      `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*SystemConfig) TableName() string {
	return TableNameSystemConfig
}

// Connection 数据库连接名称
func (m *SystemConfig) Connection() string {
	return "mysql"
}
