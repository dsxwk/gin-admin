package model

const TableNameImportRecords = "import_records"

// ImportRecords 导入记录表
type ImportRecords struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Type        int64      `gorm:"column:type;not null;default:0;type:int(10) unsigned;comment:导入类型" json:"type" form:"type"`
	Name        string     `gorm:"column:name;not null;type:varchar(20);comment:类型名称" json:"name" form:"name"`
	Data        *JsonValue `gorm:"column:data;type:json;comment:导入数据" json:"data" form:"data"`
	CreatedUser int64      `gorm:"column:created_user;not null;default:0;type:int(10) unsigned;comment:创建人" json:"createdUser" form:"createdUser"`
	User        *User      `gorm:"foreignKey:created_user;references:id" json:"user"`
	CreatedAt   *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt   *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt   *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*ImportRecords) TableName() string {
	return TableNameImportRecords
}

// Connection 数据库连接名称
func (m *ImportRecords) Connection() string {
	return "mysql"
}
