package model

const TableNameUserRoles = "user_roles"

// UserRoles 用户角色表
type UserRoles struct {
	ID        int32      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	UserId    int32      `gorm:"column:user_id;not null;default:0;type:int(10) unsigned;comment:用户id" json:"userId" form:"userId"`
	RoleId    int32      `gorm:"column:role_id;not null;default:0;type:int(10) unsigned;comment:角色id" json:"roleId" form:"roleId"`
	Name      string     `gorm:"column:name;not null;type:varchar(20);comment:角色名称" json:"name" form:"name"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*UserRoles) TableName() string {
	return TableNameUserRoles
}

// Connection 数据库连接名称
func (m *UserRoles) Connection() string {
	return "mysql"
}
