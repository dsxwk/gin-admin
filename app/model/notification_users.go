package model

const TableNameNotificationUsers = "notification_users"

// NotificationUsers 系统通知用户表
type NotificationUsers struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	NotificationId int64     `gorm:"column:notification_id;not null;default:0;type:int(10) unsigned;index:idx_notification_users_notification;comment:系统通知id" json:"notificationId" form:"notificationId"`
	ToUserId       int64     `gorm:"column:to_user_id;not null;default:0;type:int(10) unsigned;index:idx_notification_users_to_user;comment:接收用户id" json:"toUserId" form:"toUserId"`
	IsRead         int64     `gorm:"column:is_read;not null;default:2;type:tinyint(3) unsigned;index:idx_notification_users_is_read;comment:是否已读 1=是 2=否" json:"isRead" form:"isRead"`
	CreatedAt      *DateTime `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt      *DateTime `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
}

func (*NotificationUsers) TableName() string {
	return TableNameNotificationUsers
}

// Connection 数据库连接名称
func (*NotificationUsers) Connection() string {
	return "mysql"
}
