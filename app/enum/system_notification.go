package enum

import "gin/common/base"

const (
	NotificationTypeUser       int64 = 1 // 用户
	NotificationTypeDepartment int64 = 2 // 部门
	NotificationIsToAllYes     int64 = 1 // 是
	NotificationIsToAllNo      int64 = 2 // 否
	NotificationStatusNormal   int64 = 1 // 正常
	NotificationStatusRevoked  int64 = 2 // 撤回
	NotificationIsReadYes      int64 = 1 // 已读
	NotificationIsReadNo       int64 = 2 // 未读
)

// SystemNotificationEnum 系统通知枚举
type SystemNotificationEnum struct{}

// Type 通知类型
func (s *SystemNotificationEnum) Type() *base.Enum[int64] {
	return base.NewEnum(
		base.Item[int64]{Value: NotificationTypeUser, Desc: "用户"},
		base.Item[int64]{Value: NotificationTypeDepartment, Desc: "部门"},
	)
}

// IsToAll 是否推送所有
func (s *SystemNotificationEnum) IsToAll() *base.Enum[int64] {
	return base.NewEnum(
		base.Item[int64]{Value: NotificationIsToAllYes, Desc: "是"},
		base.Item[int64]{Value: NotificationIsToAllNo, Desc: "否"},
	)
}

// Status 通知状态
func (s *SystemNotificationEnum) Status() *base.Enum[int64] {
	return base.NewEnum(
		base.Item[int64]{Value: NotificationStatusNormal, Desc: "正常"},
		base.Item[int64]{Value: NotificationStatusRevoked, Desc: "撤回"},
	)
}

// IsRead 是否已读
func (s *SystemNotificationEnum) IsRead() *base.Enum[int64] {
	return base.NewEnum(
		base.Item[int64]{Value: NotificationIsReadYes, Desc: "已读"},
		base.Item[int64]{Value: NotificationIsReadNo, Desc: "未读"},
	)
}
