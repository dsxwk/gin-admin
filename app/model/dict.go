package model

import "gin/pkg"

const TableNameDict = "dict"

// Dict 字典表
type Dict struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Pid       int64          `gorm:"column:pid;not null;default:0;type:int(10) unsigned;comment:父级id" json:"pid" form:"pid"`
	Name      string         `gorm:"column:name;not null;type:varchar(50);comment:标识" json:"name" form:"name"`
	Title     string         `gorm:"column:title;not null;type:varchar(100);comment:名称" json:"title" form:"title"`
	Value     string         `gorm:"column:value;not null;type:varchar(100);comment:映射值" json:"value" form:"value"`
	Status    int64          `gorm:"column:status;not null;default:1;type:tinyint(3) unsigned;comment:状态 1=启用 2=停用" json:"status" form:"status"`
	Sort      int64          `gorm:"column:sort;not null;default:0;type:int(10) unsigned;comment:排序" json:"sort" form:"sort"`
	Extend    *JsonValue     `gorm:"column:extend;type:json;comment:扩展字段" json:"extend" form:"extend"`
	Desc      string         `gorm:"column:desc;not null;type:varchar(100);comment:字段描述" json:"desc" form:"desc"`
	Children  []pkg.TreeNode `gorm:"-;comment:子节点" json:"children"`
	CreatedAt *DateTime      `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime      `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt     `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*Dict) TableName() string {
	return TableNameDict
}

// Connection 数据库连接名称
func (m *Dict) Connection() string {
	return "mysql"
}

// GetId 实现TreeNode接口
func (m *Dict) GetId() int64 {
	return m.ID
}

func (m *Dict) GetPid() int64 {
	return m.Pid
}

func (m *Dict) GetChildren() *[]pkg.TreeNode {
	return &m.Children
}

func (m *Dict) GetTree(data []Dict) []pkg.TreeNode {
	items := make([]*Dict, 0, len(data))
	for i := range data {
		items = append(items, &data[i])
	}
	return pkg.BuildTree[*Dict](items)
}
