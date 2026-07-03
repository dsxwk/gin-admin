package model

const TableNameRoleActions = "role_actions"

// RoleActions 角色功能表
type RoleActions struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	RoleId    int64      `gorm:"column:role_id;not null;default:0;type:int(10) unsigned;comment:角色id" json:"roleId" form:"roleId"`
	ActionId  int64      `gorm:"column:action_id;not null;default:0;type:int(10) unsigned;comment:功能id" json:"actionId" form:"actionId"`
	Name      string     `gorm:"column:name;not null;type:varchar(30);comment:角色名称" json:"name" form:"name"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*RoleActions) TableName() string {
	return TableNameRoleActions
}

// Connection 数据库连接名称
func (m *RoleActions) Connection() string {
	return "mysql"
}
