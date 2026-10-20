package model

const TableNameFileGroup = "file_group"

// FileGroup 文件分组表
type FileGroup struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Pid       int64     `gorm:"column:pid;not null;default:0;type:int(10) unsigned;comment:父级ID" json:"pid" form:"pid"`
	Name      string    `gorm:"column:name;not null;type:varchar(20);comment:分组名称(如: 商品、头像)" json:"name" form:"name"`
	Sort      int64     `gorm:"column:sort;not null;default:0;type:int(10) unsigned;comment:排序(越小越前)" json:"sort" form:"sort"`
	CreatedAt *DateTime `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
}

func (*FileGroup) TableName() string {
	return TableNameFileGroup
}

// Connection 数据库连接名称
func (*FileGroup) Connection() string {
	return "mysql"
}
