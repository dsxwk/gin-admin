package model

const TableNameSystemNotification = "system_notification"

// SystemNotification 系统通知表
type SystemNotification struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	FromUserId int64     `gorm:"column:from_user_id;not null;default:0;type:int(10) unsigned;index:idx_system_notification_from_user;comment:发送用户id" json:"fromUserId" form:"fromUserId"`
	Type       int64     `gorm:"column:type;not null;default:1;type:tinyint(3) unsigned;comment:通知类型 1=用户 2=部门" json:"type" form:"type"`
	IsToAll    int64     `gorm:"column:is_to_all;not null;default:2;type:tinyint(3) unsigned;comment:是否推送所有 1=是 2=否" json:"isToAll" form:"isToAll"`
	Title      string    `gorm:"column:title;not null;type:varchar(50);comment:标题" json:"title" form:"title"`
	Content    string    `gorm:"column:content;not null;type:varchar(500);comment:内容" json:"content" form:"content"`
	Status     int64     `gorm:"column:status;not null;default:0;type:tinyint(3) unsigned;index:idx_system_notification_status;comment:状态 1=正常 2=撤回" json:"status" form:"status"`
	CreatedAt  *DateTime `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt  *DateTime `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
}

func (*SystemNotification) TableName() string {
	return TableNameSystemNotification
}

// Connection 数据库连接名称
func (*SystemNotification) Connection() string {
	return "mysql"
}
