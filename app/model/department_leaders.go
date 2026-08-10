package model

const TableNameDepartmentLeaders = "department_leaders"

// DepartmentLeaders 部门领导表
type DepartmentLeaders struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	DepartmentId int64      `gorm:"column:department_id;not null;default:0;type:int(10) unsigned;comment:部门id" json:"departmentId" form:"departmentId"`
	LeaderUserId int64      `gorm:"column:leader_user_id;not null;default:0;type:int(10) unsigned;comment:领导id" json:"leaderUserId" form:"leaderUserId"`
	Leader       *User      `gorm:"foreignKey:leader_user_id;references:id" json:"leader"`
	CreatedAt    *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt    *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt    *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*DepartmentLeaders) TableName() string {
	return TableNameDepartmentLeaders
}

// Connection 数据库连接名称
func (m *DepartmentLeaders) Connection() string {
	return "mysql"
}
