package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/app/facade"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider/ws"
)

// MessageHandler 消息处理
type MessageHandler struct{}

type MessageRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type bindMessage struct {
	UserID int64 `json:"userId"`
}

// Handle 处理客户端消息
func (h *MessageHandler) Handle(_ context.Context, client ws.Connection, message ws.Message) (err error) {
	facade.Log().Debug(fmt.Sprintf("ws收到消息: clientId=%s userId=%d data=%s", client.ID(), client.UserID(), string(message.Data)))
	defer func() {
		if err != nil {
			facade.Log().Warn(fmt.Sprintf("ws消息处理失败: clientId=%s userId=%d err=%v", client.ID(), client.UserID(), err))
		}
	}()

	var request MessageRequest
	if err = json.Unmarshal(message.Data, &request); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("消息格式错误"))
	}

	switch request.Type {
	case "bind":
		var payload bindMessage
		if err = json.Unmarshal(request.Data, &payload); err != nil {
			return ws.SendError(client, errcode.ArgsError().WithMsg("userId参数错误"))
		}
		if payload.UserID <= 0 {
			return ws.SendError(client, errcode.ArgsError().WithMsg("userId不能为空"))
		}
		if boundUserID := client.UserID(); boundUserID > 0 && boundUserID != payload.UserID {
			return ws.SendError(client, errcode.Forbidden().WithMsg("已绑定登录用户,不能切换用户"))
		}

		client.BindUser(payload.UserID)
		facade.Log().Info(fmt.Sprintf("ws绑定用户成功: clientId=%s userId=%d", client.ID(), payload.UserID))
		return ws.SendJSON(client, "bound", map[string]any{
			"clientId": client.ID(),
			"userId":   payload.UserID,
		})

	case "ping":
		return ws.SendJSON(client, "pong", nil)

	case "chat":
		var payload any
		if len(request.Data) > 0 {
			if err = json.Unmarshal(request.Data, &payload); err != nil {
				return ws.SendError(client, errcode.ArgsError().WithMsg("消息内容格式错误"))
			}
		}
		facade.Log().Info(fmt.Sprintf("ws收到聊天消息: clientId=%s userId=%d", client.ID(), client.UserID()))
		return ws.SendJSON(client, "chat", payload)
	}

	facade.Log().Debug(fmt.Sprintf("ws不支持的消息类型: clientId=%s type=%s", client.ID(), request.Type))
	return ws.SendError(client, errcode.ArgsError().WithMsg("不支持的消息类型: "+request.Type))
}
