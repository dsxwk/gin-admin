package model

import "gin/pkg"

const TableNameMenuActions = "menu_actions"

// MenuActions 菜单功能表
type MenuActions struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Pid         int64          `gorm:"column:pid;not null;default:0;type:int(10) unsigned;comment:父级id" json:"pid" form:"pid"`
	MenuId      int64          `gorm:"column:menu_id;not null;default:0;type:int(10) unsigned;comment:菜单id" json:"menuId" form:"menuId"`
	Type        int64          `gorm:"column:type;not null;default:1;type:tinyint(3) unsigned;comment:类型 1=header 2=operation" json:"type" form:"type"`
	BtnType     string         `gorm:"column:btn_type;not null;default:btn;type:varchar(20);comment:按钮类型 text|btn" json:"btnType" form:"btnType"`
	BtnStyle    string         `gorm:"column:btn_style;not null;default:primary;type:varchar(20);comment:按钮样式" json:"btnStyle" form:"btnStyle"`
	BtnSize     string         `gorm:"column:btn_size;not null;default:small;type:varchar(20);comment:按钮尺寸" json:"btnSize" form:"btnSize"`
	IsConfirm   int64          `gorm:"column:is_confirm;not null;default:2;type:tinyint(3) unsigned;comment:是否确认 1=是 2=否" json:"isConfirm" form:"isConfirm"`
	Label       string         `gorm:"column:label;not null;type:varchar(30);comment:功能名称" json:"label" form:"label"`
	AuthValue   string         `gorm:"column:auth_value;not null;type:varchar(100);comment:权限标识" json:"authValue" form:"authValue"`
	IsLink      int64          `gorm:"column:is_link;not null;default:2;type:tinyint(3) unsigned;comment:是否为链接 1=是 2=否" json:"isLink" form:"isLink"`
	Sort        int64          `gorm:"column:sort;not null;default:0;type:int(10) unsigned;comment:排序" json:"sort" form:"sort"`
	RoleActions []*RoleActions `gorm:"foreignKey:action_id;references:id;comment:角色功能" json:"roleActions"`
	Parent      *MenuActions   `gorm:"foreignKey:pid;references:id;comment:父级功能" json:"parent"`
	Children    []pkg.TreeNode `gorm:"-;comment:子节点" json:"children"`
	CreatedAt   *DateTime      `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt   *DateTime      `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt   *DeletedAt     `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*MenuActions) TableName() string {
	return TableNameMenuActions
}

// Connection 数据库连接名称
func (m *MenuActions) Connection() string {
	return "mysql"
}

// GetId 实现TreeNode接口
func (m *MenuActions) GetId() int64 {
	return m.ID
}

func (m *MenuActions) GetPid() int64 {
	return m.Pid
}

func (m *MenuActions) GetChildren() *[]pkg.TreeNode {
	return &m.Children
}

func (m *MenuActions) GetTree(data []MenuActions) []pkg.TreeNode {
	items := make([]*MenuActions, 0, len(data))
	for i := range data {
		items = append(items, &data[i])
	}
	return pkg.BuildTree[*MenuActions](items)
}
