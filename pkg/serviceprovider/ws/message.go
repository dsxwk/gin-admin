package ws

import (
	"encoding/json"
	"errors"
	"gin/pkg/errcode"

	"github.com/gorilla/websocket"
)

// MessageType WebSocket消息类型
type MessageType int

const (
	// MessageText 文本消息
	MessageText MessageType = websocket.TextMessage
	// MessageBinary 二进制消息
	MessageBinary MessageType = websocket.BinaryMessage
)

// Message WebSocket消息
type Message struct {
	Type MessageType // 消息类型
	Data []byte      // 消息内容
}

// Envelope WebSocket JSON消息信封
type Envelope struct {
	Type string `json:"type"`           // 消息类型
	Data any    `json:"data,omitempty"` // 消息数据
}

// errorMessage WebSocket错误消息
type errorMessage struct {
	Code    int64  `json:"code"`    // 错误码
	Message string `json:"message"` // 错误描述
}

// SendJSON 发送JSON消息
func SendJSON(connection Connection, messageType string, data any) error {
	if connection == nil {
		return ErrClosed
	}

	payload, err := json.Marshal(Envelope{
		Type: messageType,
		Data: data,
	})
	if err != nil {
		return err
	}

	return connection.Send(Message{
		Type: MessageText,
		Data: payload,
	})
}

// SendError 发送公共错误消息
func SendError(connection Connection, err error) error {
	if err == nil {
		return nil
	}

	if errorCode, ok := errors.AsType[errcode.ErrorCode](err); ok {
		return SendJSON(connection, "error", errorMessage{
			Code:    errorCode.Code,
			Message: errorCode.Msg,
		})
	}

	return SendJSON(connection, "error", errorMessage{
		Code:    500,
		Message: err.Error(),
	})
}
