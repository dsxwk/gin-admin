package request

import (
	"gin/app/enum"
	"gin/app/errcode"
	"gin/common/base"
	"strings"

	"github.com/gookit/validate"
)

// SystemNotificationPush 推送通知请求
type SystemNotificationPush struct {
	base.BaseRequest
	Title         string  `json:"title" validate:"required|minLen:1|maxLen:50" label:"通知标题"`    // 标题
	Content       string  `json:"content" validate:"required|minLen:1|maxLen:500" label:"通知内容"` // 内容
	Type          int64   `json:"type" validate:"int|in:1,2" label:"通知类型 1=用户 2=部门"`            // 通知类型 1=用户 2=部门
	IsToAll       int64   `json:"isToAll" validate:"int|in:1,2" label:"是否推送所有 1=是 2=否"`         // 是否推送所有 1=是 2=否
	UserIds       []int64 `json:"userIds" validate:"" label:"接收用户ID"`                           // 接收用户ID
	DepartmentIds []int64 `json:"departmentIds" validate:"" label:"接收部门ID"`                     // 接收部门ID
}

// SystemNotificationList 通知列表请求
type SystemNotificationList struct {
	base.BaseRequest
	Page     int  `json:"page" validate:"int|gte:1" label:"页码"`       // 页码
	PageSize int  `json:"pageSize" validate:"int|gte:1" label:"每页数量"` // 每页数量
	Sent     bool `json:"sent" label:"已发送列表"`                         // 是否查询已发送列表
}

// SystemNotificationRead 通知已读请求
type SystemNotificationRead struct {
	base.BaseRequest
	IDs []int64 `json:"ids" validate:"" label:"通知ID列表"` // 通知ID列表,为空表示全部已读
}

// SystemNotificationRevoke 通知撤回请求
type SystemNotificationRevoke struct {
	base.BaseRequest
	ID int64 `json:"id" validate:"required|int|gt:0" label:"通知ID"` // 通知ID
}

// Validate 推送通知验证
func (s *SystemNotificationPush) Validate(scene string) error {
	s.Title = strings.TrimSpace(s.Title)
	s.Content = strings.TrimSpace(s.Content)
	if s.Type == 0 {
		s.Type = enum.NotificationTypeUser
	}
	if s.IsToAll == 0 {
		s.IsToAll = enum.NotificationIsToAllNo
	}

	if err := validateSystemNotification(s, scene); err != nil {
		return err
	}
	if s.IsToAll == enum.NotificationIsToAllYes {
		return nil
	}

	switch s.Type {
	case enum.NotificationTypeUser:
		if !validPositiveSystemNotificationIDs(s.UserIds) {
			return errcode.ArgsError().WithMsg("接收用户不能为空")
		}
	case enum.NotificationTypeDepartment:
		if !validPositiveSystemNotificationIDs(s.DepartmentIds) {
			return errcode.ArgsError().WithMsg("接收部门不能为空")
		}
	}
	return nil
}

// Validate 通知列表验证
func (s *SystemNotificationList) Validate(scene string) error {
	if s.Page == 0 {
		s.Page = 1
	}
	if s.PageSize == 0 {
		s.PageSize = 20
	}
	return validateSystemNotification(s, scene)
}

// Validate 通知已读验证
func (s *SystemNotificationRead) Validate(scene string) error {
	for _, id := range s.IDs {
		if id <= 0 {
			return errcode.ArgsError().WithMsg("通知ID必须大于0")
		}
	}
	return validateSystemNotification(s, scene)
}

// Validate 通知撤回验证
func (s *SystemNotificationRevoke) Validate(scene string) error {
	return validateSystemNotification(s, scene)
}

// ConfigValidation 推送通知验证场景
func (s *SystemNotificationPush) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Push": []string{"Title", "Content", "Type", "IsToAll", "UserIds", "DepartmentIds"},
	})
}

// ConfigValidation 通知列表验证场景
func (s *SystemNotificationList) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"List": []string{"Page", "PageSize", "Sent"},
	})
}

// ConfigValidation 通知已读验证场景
func (s *SystemNotificationRead) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Read": []string{"IDs"},
	})
}

// ConfigValidation 通知撤回验证场景
func (s *SystemNotificationRevoke) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Revoke": []string{"ID"},
	})
}

// Messages 验证器错误消息
func (s *SystemNotificationPush) Messages() map[string]string {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationList) Messages() map[string]string {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationRead) Messages() map[string]string {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationRevoke) Messages() map[string]string {
	return systemNotificationMessages()
}

// Translates 字段翻译
func (s *SystemNotificationPush) Translates() map[string]string {
	return validate.MS{
		"Title":         "通知标题",
		"Content":       "通知内容",
		"Type":          "通知类型",
		"IsToAll":       "是否推送所有",
		"UserIds":       "接收用户ID",
		"DepartmentIds": "接收部门ID",
	}
}

// Translates 字段翻译
func (s *SystemNotificationList) Translates() map[string]string {
	return validate.MS{
		"Page":     "页码",
		"PageSize": "每页数量",
		"Sent":     "已发送列表",
	}
}

// Translates 字段翻译
func (s *SystemNotificationRead) Translates() map[string]string {
	return validate.MS{
		"IDs": "通知ID列表",
	}
}

// Translates 字段翻译
func (s *SystemNotificationRevoke) Translates() map[string]string {
	return validate.MS{
		"ID": "通知ID",
	}
}

// validateSystemNotification 验证系统通知请求
func validateSystemNotification(data any, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// systemNotificationMessages 系统通知验证消息
func systemNotificationMessages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"in":       "字段 {field} 值不正确",
		"minLen":   "{field} 长度不能少于 {min} 个字符",
		"maxLen":   "{field} 长度不能超过 {max} 个字符",
		"gt":       "字段 {field} 必须大于 0",
		"gte":      "字段 {field} 必须大于等于 {min}",
	}
}

// validPositiveSystemNotificationIDs 是否全部为有效ID
func validPositiveSystemNotificationIDs(ids []int64) bool {
	if len(ids) == 0 {
		return false
	}
	for _, id := range ids {
		if id <= 0 {
			return false
		}
	}
	return true
}
