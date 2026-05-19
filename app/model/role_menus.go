package model

const TableNameRoleMenus = "role_menus"

type RoleMenus struct {
	Id        int64      `gorm:"column:id;comment:ID" json:"id" form:"id"`
	RoleId    int64      `gorm:"column:role_id;comment:角色id" json:"roleId" form:"roleId"`
	MenuId    int64      `gorm:"column:menu_id;comment:菜单id" json:"menuId" form:"menuId"`
	Name      string     `gorm:"column:name;comment:角色名称" json:"name" form:"name"`
	CreatedAt *DateTime  `gorm:"column:created_at;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*RoleMenus) TableName() string {
	return TableNameRoleMenus
}

// Connection 数据库连接名称
func (m *RoleMenus) Connection() string {
	return "mysql"
}
