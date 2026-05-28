package model

const TableNameRoleMenus = "role_menus"

// RoleMenus 角色菜单表
type RoleMenus struct {
	Id        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	RoleId    int64      `gorm:"column:role_id;not null;default:0;type:int(10) unsigned;comment:角色id" json:"roleId" form:"roleId"`
	MenuId    int64      `gorm:"column:menu_id;not null;default:0;type:int(10) unsigned;comment:菜单id" json:"menuId" form:"menuId"`
	Name      string     `gorm:"column:name;not null;type:varchar(20);comment:角色名称" json:"name" form:"name"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*RoleMenus) TableName() string {
	return TableNameRoleMenus
}

// Connection 数据库连接名称
func (m *RoleMenus) Connection() string {
	return "mysql"
}
