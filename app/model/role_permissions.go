package model

const TableNameRolePermissions = "role_permissions"

// RolePermissions 角色权限表
type RolePermissions struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	RoleId       int64      `gorm:"column:role_id;not null;default:0;type:int(10) unsigned;comment:角色id" json:"roleId" form:"roleId"`
	PermissionId int64      `gorm:"column:permission_id;not null;default:0;type:int(10) unsigned;comment:权限id" json:"permissionId" form:"permissionId"`
	UpdatedAt    *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt    *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*RolePermissions) TableName() string {
	return TableNameRolePermissions
}

// Connection 数据库连接名称
func (m *RolePermissions) Connection() string {
	return "mysql"
}
