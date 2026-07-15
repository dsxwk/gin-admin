package model

const TableNameMenuMeta = "menu_meta"

// MenuMeta 菜单元数据表
type MenuMeta struct {
	ID          int64          `gorm:"column:id;type:int(10) unsigned;primaryKey;autoIncrement:true;comment:ID" json:"id"`
	MenuId      int64          `gorm:"column:menu_id;type:int(10) unsigned;comment:菜单ID" json:"menuId"`
	Title       string         `gorm:"column:title;type:json;comment:菜单名称" json:"title"`
	Icon        string         `gorm:"column:icon;type:json;comment:菜单图标" json:"icon"`
	Path        string         `gorm:"column:path;type:varchar(50);not null;comment:路由路径" json:"path"`        // 路由路径
	Redirect    string         `gorm:"column:redirect;type:varchar(50);not null;comment:重定向" json:"redirect"` // 重定向
	Component   string         `gorm:"column:component;type:varchar(100);not null;comment:组件路径" json:"component"`
	IsHide      int64          `gorm:"column:is_hide;type:json;comment:是否隐藏 1=是 2=否" json:"isHide"`
	IsKeepAlive int64          `gorm:"column:is_keep_alive;type:json;comment:是否缓存 1=是 2=否" json:"isKeepAlive"`
	IsAffix     int64          `gorm:"column:is_affix;type:json;comment:是否固定 1=是 2=否" json:"isAffix"`
	IsLink      string         `gorm:"column:is_link;type:json;comment:外链/内嵌时链接地址" json:"isLink"`               // 外链/内嵌时链接地址(http:xxx.com),开启外链条件`1 isLink:链接地址不为空`
	IsIframe    int64          `gorm:"column:is_iframe;type:json;comment:是否内嵌 1=是 2=否" json:"isIframe"`         // 是否内嵌,开启条件`1 isIframe:true 2 isLink:链接地址不为空`
	Roles       []*RoleMenus   `gorm:"foreignKey:menu_id;references:menu_id;comment:菜单角色" json:"roles"`         // 权限标识,取角色管理
	AuthBtnList []*MenuActions `gorm:"foreignKey:menu_id;references:menu_id;comment:按钮权限列表" json:"authBtnList"` // 按钮权限列表
	CreatedAt   *DateTime      `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt"`
	UpdatedAt   *DateTime      `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *DeletedAt     `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" swaggerignore:"true"`
}

func (*MenuMeta) TableName() string {
	return TableNameMenuMeta
}

// Connection 数据库连接名称
func (m *MenuMeta) Connection() string {
	return "mysql"
}
