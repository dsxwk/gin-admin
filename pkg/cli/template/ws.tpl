package {{.Package}}

import (
	"context"
	"encoding/json"
	"gin/pkg/serviceprovider/ws"
)

const {{.Name}}MessageType = "{{.MessageType}}" // {{.Description}}

// {{.Name}}Handler {{.Description}}
type {{.Name}}Handler struct{}

// {{.Name}}Request {{.Description}}请求
type {{.Name}}Request struct {
	Type string          `json:"type"` // 消息类型
	Data json.RawMessage `json:"data"` // 消息数据
}

// Match 匹配{{.Description}}
func (h *{{.Name}}Handler) Match(message ws.Message) bool {
	if message.Type != ws.MessageText {
		return false
	}

	var request {{.Name}}Request
	if err := json.Unmarshal(message.Data, &request); err != nil {
		return false
	}

	return request.Type == {{.Name}}MessageType
}

// Handle {{.Description}}
func (h *{{.Name}}Handler) Handle(_ context.Context, _ ws.Connection, _ ws.Message) error {
	// TODO:处理业务逻辑
	return nil
}
