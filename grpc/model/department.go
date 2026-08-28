package model

const TableNameDepartment = "department"

// Department 部门表
type Department struct {
	ID        int32      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Pid       int32      `gorm:"column:pid;not null;default:0;type:int(10) unsigned;comment:父级id" json:"pid" form:"pid"`
	Name      string     `gorm:"column:name;not null;type:varchar(50);comment:部门名称" json:"name" form:"name"`
	Status    int32      `gorm:"column:status;not null;default:1;type:tinyint(3) unsigned;comment:状态 1=启用 2=停用" json:"status" form:"status"`
	Sort      int32      `gorm:"column:sort;not null;default:0;type:int(10) unsigned;comment:排序" json:"sort" form:"sort"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*Department) TableName() string {
	return TableNameDepartment
}

// Connection 数据库连接名称
func (m *Department) Connection() string {
	return "mysql"
}
