package model

const TableNameMenuAction = "menu_actions"

type MenuActions struct {
	Id          int64          `gorm:"column:id;comment:ID" json:"id" form:"id"`
	Pid         int64          `gorm:"column:pid;comment:父级id" json:"pid" form:"pid"`
	MenuId      int64          `gorm:"column:menu_id;comment:菜单id" json:"menuId" form:"menuId"`
	Type        int64          `gorm:"column:type;comment:类型 1=header 2=operation" json:"type" form:"type"`
	BtnType     string         `gorm:"column:btn_type;comment:按钮类型 text|btn" json:"btnType" form:"btnType"`
	BtnStyle    string         `gorm:"column:btn_style;comment:按钮样式" json:"btnStyle" form:"btnStyle"`
	BtnSize     string         `gorm:"column:btn_size;comment:按钮尺寸" json:"btnSize" form:"btnSize"`
	IsConfirm   int64          `gorm:"column:is_confirm;comment:是否确认 1=是 2=否" json:"isConfirm" form:"isConfirm"`
	Label       string         `gorm:"column:label;comment:功能名称" json:"label" form:"label"`
	AuthValue   string         `gorm:"column:auth_value;comment:权限标识" json:"authValue" form:"authValue"`
	IsLink      int64          `gorm:"column:is_link;comment:是否为链接 1=是 2=否" json:"isLink" form:"isLink"`
	Sort        int64          `gorm:"column:sort;comment:排序" json:"sort" form:"sort"`
	RoleActions []*RoleActions `gorm:"foreignKey:action_id;references:id;comment:角色功能" json:"roleActions"`
	Parent      *MenuActions   `gorm:"foreignKey:pid;references:id;comment:父级功能" json:"parent"`
	CreatedAt   *DateTime      `gorm:"column:created_at;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt   *DateTime      `gorm:"column:updated_at;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt   *DeletedAt     `gorm:"column:deleted_at;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*MenuActions) TableName() string {
	return TableNameMenuAction
}

// Connection 数据库连接名称
func (m *MenuActions) Connection() string {
	return "mysql"
}
