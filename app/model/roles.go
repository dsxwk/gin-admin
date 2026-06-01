package model

const TableNameRoles = "roles"

// Roles 角色表
type Roles struct {
	Id        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Name      string     `gorm:"column:name;not null;type:varchar(20);comment:角色名称" json:"name" form:"name"`
	Desc      string     `gorm:"column:desc;not null;type:varchar(100);comment:角色描述" json:"desc" form:"desc"`
	Status    int64      `gorm:"column:status;not null;default:1;type:tinyint(1) unsigned;comment:状态 1=启用 2=停用" json:"status" form:"status"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*Roles) TableName() string {
	return TableNameRoles
}

// Connection 数据库连接名称
func (m *Roles) Connection() string {
	return "mysql"
}
