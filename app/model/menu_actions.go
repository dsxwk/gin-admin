package model

const TableNameMenuActions = "menu_actions"

// MenuActions 菜单功能表
type MenuActions struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	MenuId    int64      `gorm:"column:menu_id;not null;default:0;type:int(10) unsigned;comment:菜单id" json:"menuId" form:"menuId"`
	Type      int64      `gorm:"column:type;not null;default:1;type:tinyint(3) unsigned;comment:类型 1=header 2=operation" json:"type" form:"type"`
	BtnType   string     `gorm:"column:btn_type;not null;default:btn;type:varchar(20);comment:按钮类型 text|btn" json:"btnType" form:"btnType"`
	BtnStyle  string     `gorm:"column:btn_style;not null;default:primary;type:varchar(20);comment:按钮样式" json:"btnStyle" form:"btnStyle"`
	BtnSize   string     `gorm:"column:btn_size;not null;default:small;type:varchar(20);comment:按钮尺寸" json:"btnSize" form:"btnSize"`
	IsConfirm int64      `gorm:"column:is_confirm;not null;default:2;type:tinyint(3) unsigned;comment:是否确认 1=是 2=否" json:"isConfirm" form:"isConfirm"`
	Label     string     `gorm:"column:label;not null;type:varchar(30);comment:功能名称" json:"label" form:"label"`
	AuthValue string     `gorm:"column:auth_value;not null;type:varchar(100);comment:权限标识" json:"authValue" form:"authValue"`
	IsLink    int64      `gorm:"column:is_link;not null;default:2;type:tinyint(3) unsigned;comment:是否为链接 1=是 2=否" json:"isLink" form:"isLink"`
	CreatedAt *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*MenuActions) TableName() string {
	return TableNameMenuActions
}

// Connection 数据库连接名称
func (m *MenuActions) Connection() string {
	return "mysql"
}
