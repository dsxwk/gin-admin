package request

import (
	"gin/app/enum"
	"gin/common/base"
	"strings"

	"github.com/dsxwk/validatex"
)

// SystemNotificationPush 推送通知请求
type SystemNotificationPush struct {
	base.Request
	Title         string  `json:"title" form:"title" label:"通知标题"`
	Content       string  `json:"content" form:"content" label:"通知内容"`
	Type          int64   `json:"type" form:"type" label:"通知类型 1=用户 2=部门"`
	IsToAll       int64   `json:"isToAll" form:"isToAll" label:"是否推送所有 1=是 2=否"`
	UserIds       []int64 `json:"userIds" form:"userIds" label:"接收用户ID"`
	DepartmentIds []int64 `json:"departmentIds" form:"departmentIds" label:"接收部门ID"`
}

// SystemNotificationList 通知列表请求
type SystemNotificationList struct {
	base.Request
	Page     int  `json:"page" form:"page" label:"页码"`
	PageSize int  `json:"pageSize" form:"pageSize" label:"每页数量"`
	Sent     bool `json:"sent" form:"sent" label:"已发送列表"`
}

// SystemNotificationRead 通知已读请求
type SystemNotificationRead struct {
	base.Request
	IDs []int64 `json:"ids" form:"ids" label:"通知ID列表"`
}

// SystemNotificationRevoke 通知撤回请求
type SystemNotificationRevoke struct {
	base.Request
	ID int64 `json:"id" form:"id" label:"通知ID"`
}

// ValidateWithScene 推送通知验证
func (s *SystemNotificationPush) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Push":
		s.Title = strings.TrimSpace(s.Title)
		s.Content = strings.TrimSpace(s.Content)
		if s.Type == 0 {
			s.Type = enum.NotificationTypeUser
		}
		if s.IsToAll == 0 {
			s.IsToAll = enum.NotificationIsToAllNo
		}

		errs := validatex.Check(
			fields.Str("Title", s.Title).Required().Max(50).Errors(),
			fields.Str("Content", s.Content).Required().Max(500).Errors(),
			fields.Int("Type", s.Type).In(enum.NotificationTypeUser, enum.NotificationTypeDepartment).Errors(),
			fields.Int("IsToAll", s.IsToAll).In(enum.NotificationIsToAllYes, enum.NotificationIsToAllNo).Errors(),
		)
		if len(errs) > 0 || s.IsToAll == enum.NotificationIsToAllYes {
			return errs
		}

		switch s.Type {
		case enum.NotificationTypeUser:
			return fields.Slice("UserIds", s.UserIds).Required().Each("gt", "", func(id int64) bool {
				return id > 0
			}, int64(0)).Errors()
		case enum.NotificationTypeDepartment:
			return fields.Slice("DepartmentIds", s.DepartmentIds).Required().Each("gt", "", func(id int64) bool {
				return id > 0
			}, int64(0)).Errors()
		default:
			return errs
		}
	default:
		return nil
	}
}

// ValidateWithScene 通知列表验证
func (s *SystemNotificationList) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		if s.Page == 0 {
			s.Page = 1
		}
		if s.PageSize == 0 {
			s.PageSize = 20
		}
		return validatex.Check(
			fields.Int("Page", s.Page).Gte(1).Errors(),
			fields.Int("PageSize", s.PageSize).Gte(1).Errors(),
		)
	default:
		return nil
	}
}

// ValidateWithScene 通知已读验证
func (s *SystemNotificationRead) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Read":
		if len(s.IDs) == 0 {
			return nil
		}
		return fields.Slice("IDs", s.IDs).Each("gt", "", func(id int64) bool {
			return id > 0
		}, int64(0)).Errors()
	default:
		return nil
	}
}

// ValidateWithScene 通知撤回验证
func (s *SystemNotificationRevoke) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Revoke":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// systemNotificationMessages 系统通知验证消息
func systemNotificationMessages() validatex.Messages {
	messages := commonMessages()
	messages["in"] = "字段 {field} 值不正确"
	messages["gte"] = "{field} 必须大于等于 {min}"
	return messages
}

// Messages 验证器错误消息
func (s *SystemNotificationPush) Messages() validatex.Messages {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationList) Messages() validatex.Messages {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationRead) Messages() validatex.Messages {
	return systemNotificationMessages()
}

// Messages 验证器错误消息
func (s *SystemNotificationRevoke) Messages() validatex.Messages {
	return systemNotificationMessages()
}

// Translates 字段翻译
func (s *SystemNotificationPush) Translates() validatex.Translates {
	return validatex.Translates{
		"Title":         "通知标题",
		"Content":       "通知内容",
		"Type":          "通知类型",
		"IsToAll":       "是否推送所有",
		"UserIds":       "接收用户ID",
		"DepartmentIds": "接收部门ID",
	}
}

// Translates 字段翻译
func (s *SystemNotificationList) Translates() validatex.Translates {
	return validatex.Translates{
		"Page":     "页码",
		"PageSize": "每页数量",
		"Sent":     "已发送列表",
	}
}

// Translates 字段翻译
func (s *SystemNotificationRead) Translates() validatex.Translates {
	return validatex.Translates{"IDs": "通知ID列表"}
}

// Translates 字段翻译
func (s *SystemNotificationRevoke) Translates() validatex.Translates {
	return validatex.Translates{"ID": "通知ID"}
}
