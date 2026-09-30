package websocket

import (
	"context"

	"gin/app/enum"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider/ws"

	"github.com/goccy/go-json"
)

const (
	notificationPushMessage    = "notification.push"
	notificationListMessage    = "notification.list"
	notificationReadMessage    = "notification.read"
	notificationRevokeMessage  = "notification.revoke"
	notificationCreatedMessage = "notification.created"
	notificationRevokedMessage = "notification.revoked"
)

// NotificationHandler 系统通知消息处理
type NotificationHandler struct{}

// Match 匹配通知消息
func (h *NotificationHandler) Match(message ws.Message) bool {
	var requestData MessageRequest
	if err := json.Unmarshal(message.Data, &requestData); err != nil {
		return false
	}

	switch requestData.Type {
	case notificationPushMessage, notificationListMessage, notificationReadMessage, notificationRevokeMessage:
		return true
	default:
		return false
	}
}

// Handle 处理通知消息
func (h *NotificationHandler) Handle(ctx context.Context, client ws.Connection, message ws.Message) error {
	userID := ws.UserID(ctx)
	if userID <= 0 {
		return ws.SendError(client, errcode.Unauthorized().WithMsg("请先登录"))
	}

	var requestData MessageRequest
	if err := json.Unmarshal(message.Data, &requestData); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("消息格式错误"))
	}

	switch requestData.Type {
	case notificationPushMessage:
		return h.push(ctx, client, userID, requestData.Data)
	case notificationListMessage:
		return h.list(ctx, client, userID, requestData.Data)
	case notificationReadMessage:
		return h.read(ctx, client, userID, requestData.Data)
	case notificationRevokeMessage:
		return h.revoke(ctx, client, userID, requestData.Data)
	default:
		return ws.SendError(client, errcode.ArgsError().WithMsg("不支持的通知消息类型"))
	}
}

// push 推送通知
func (h *NotificationHandler) push(ctx context.Context, client ws.Connection, userID int64, data json.RawMessage) error {
	var payload request.SystemNotificationPush
	if err := json.Unmarshal(data, &payload); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("推送参数错误"))
	}
	if err := payload.Validate("Push"); err != nil {
		return ws.SendError(client, err)
	}

	var notificationService service.SystemNotificationService
	notification, userIDs, err := notificationService.Push(ctx, userID, payload)
	if err != nil {
		return ws.SendError(client, err)
	}

	message := notificationPayload{
		ID:         notification.ID,
		FromUserId: notification.FromUserId,
		Type:       notification.Type,
		IsToAll:    notification.IsToAll,
		Title:      notification.Title,
		Content:    notification.Content,
		Status:     notification.Status,
		IsRead:     enum.NotificationIsReadNo,
		CreatedAt:  notification.CreatedAt,
	}
	sent := sendNotificationToUsers(userIDs, notificationCreatedMessage, message)

	return ws.SendJSON(client, notificationPushMessage, map[string]any{
		"id":         notification.ID,
		"recipients": len(userIDs),
		"sent":       sent,
	})
}

// list 获取通知列表
func (h *NotificationHandler) list(ctx context.Context, client ws.Connection, userID int64, data json.RawMessage) error {
	var payload request.SystemNotificationList
	if len(data) > 0 {
		if err := json.Unmarshal(data, &payload); err != nil {
			return ws.SendError(client, errcode.ArgsError().WithMsg("列表参数错误"))
		}
	}
	if err := payload.Validate("List"); err != nil {
		return ws.SendError(client, err)
	}

	var notificationService service.SystemNotificationService
	list, err := notificationService.List(ctx, userID, payload.Page, payload.PageSize, payload.Sent)
	if err != nil {
		return ws.SendError(client, err)
	}

	return ws.SendJSON(client, notificationListMessage, list)
}

// read 标记通知已读
func (h *NotificationHandler) read(ctx context.Context, client ws.Connection, userID int64, data json.RawMessage) error {
	var payload request.SystemNotificationRead
	if len(data) > 0 {
		if err := json.Unmarshal(data, &payload); err != nil {
			return ws.SendError(client, errcode.ArgsError().WithMsg("已读参数错误"))
		}
	}
	if err := payload.Validate("Read"); err != nil {
		return ws.SendError(client, err)
	}

	var notificationService service.SystemNotificationService
	updated, err := notificationService.Read(ctx, userID, payload.IDs)
	if err != nil {
		return ws.SendError(client, err)
	}

	return ws.SendJSON(client, notificationReadMessage, map[string]any{
		"updated": updated,
	})
}

// revoke 撤回通知
func (h *NotificationHandler) revoke(ctx context.Context, client ws.Connection, userID int64, data json.RawMessage) error {
	var payload request.SystemNotificationRevoke
	if err := json.Unmarshal(data, &payload); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("撤回参数错误"))
	}
	if err := payload.Validate("Revoke"); err != nil {
		return ws.SendError(client, err)
	}

	var notificationService service.SystemNotificationService
	notification, userIDs, err := notificationService.Revoke(ctx, userID, payload.ID)
	if err != nil {
		return ws.SendError(client, err)
	}

	sendNotificationToUsers(userIDs, notificationRevokedMessage, map[string]any{
		"id":     notification.ID,
		"status": notification.Status,
	})

	return ws.SendJSON(client, notificationRevokeMessage, map[string]any{
		"id":     notification.ID,
		"status": notification.Status,
	})
}

// notificationPayload 通知推送数据
type notificationPayload struct {
	ID         int64  `json:"id"`         // 通知ID
	FromUserId int64  `json:"fromUserId"` // 发送用户ID
	Type       int64  `json:"type"`       // 通知类型 1=用户 2=部门
	IsToAll    int64  `json:"isToAll"`    // 是否推送所有 1=是 2=否
	Title      string `json:"title"`      // 标题
	Content    string `json:"content"`    // 内容
	Status     int64  `json:"status"`     // 状态 1=正常 2=撤回
	IsRead     int64  `json:"isRead"`     // 是否已读 1=是 2=否
	CreatedAt  any    `json:"createdAt"`  // 创建时间
}

// sendNotificationToUsers 向用户推送通知消息
func sendNotificationToUsers(userIDs []int64, messageType string, data any) int {
	manager := facade.WS()
	if manager == nil {
		return 0
	}

	payload, err := json.Marshal(ws.Envelope{
		Type: messageType,
		Data: data,
	})
	if err != nil {
		return 0
	}

	message := ws.Message{
		Type: ws.MessageText,
		Data: payload,
	}
	return manager.SendToUsers(userIDs, message)
}
