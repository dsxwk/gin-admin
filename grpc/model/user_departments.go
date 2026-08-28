package model

const TableNameUserDepartments = "user_departments"

// UserDepartments 用户部门表
type UserDepartments struct {
	ID           int32       `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	UserId       int32       `gorm:"column:user_id;not null;default:0;type:int(10) unsigned;comment:用户id" json:"userId" form:"userId"`
	DepartmentId int32       `gorm:"column:department_id;not null;default:0;type:int(10) unsigned;comment:部门id" json:"departmentId" form:"departmentId"`
	IsMain       int32       `gorm:"column:is_main;not null;default:2;type:tinyint(3) unsigned;comment:是否是主部门 1=是 2=否" json:"isMain" form:"isMain"`
	Dept         *Department `gorm:"foreignKey:department_id;references:id;comment:部门" json:"dept"`
	CreatedAt    *DateTime   `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt    *DateTime   `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt    *DeletedAt  `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*UserDepartments) TableName() string {
	return TableNameUserDepartments
}

// Connection 数据库连接名称
func (m *UserDepartments) Connection() string {
	return "mysql"
}
