package service

import (
	"gin/app/model"
	"gin/common/base"
	"gin/pkg/serviceprovider/agent"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

// AgentService AI会话服务
type AgentService struct {
	base.BaseService
}

// CreateSession 创建会话
func (s *AgentService) CreateSession(userId int64, title, provider, modelName, traceId string) (int64, error) {
	session := model.AgentSession{
		UserId:   userId,
		Title:    title,
		Provider: provider,
		Model:    modelName,
		TraceId:  traceId,
	}

	db := s.DB(&session)
	if err := db.Model(&session).Create(&session).Error; err != nil {
		return 0, err
	}
	return session.ID, nil
}

// GetSession 获取会话
func (s *AgentService) GetSession(id int64) (session model.AgentSession, err error) {
	db := s.DB(&session)
	err = db.Model(&session).First(&session, id).Error
	return session, err
}

// ListSessions 获取用户会话列表
func (s *AgentService) ListSessions(userId int64) (sessions []model.AgentSession, err error) {
	var session model.AgentSession
	db := s.DB(&session)
	err = db.Model(&session).Where("user_id = ?", userId).Order("id DESC").Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// RecordMessage 记录消息(实现agent.Recorder接口)
func (s *AgentService) RecordMessage(sessionId int64, record agent.MessageRecord) error {
	var (
		m  model.AgentMessage
		db = s.DB(&m)
	)

	msg := model.AgentMessage{
		SessionId:  sessionId,
		Role:       record.Role,
		Content:    record.Content,
		ToolName:   record.ToolName,
		ToolCallId: record.ToolCallId,
		Tokens:     record.Tokens,
		CostMs:     record.CostMs,
	}

	if record.ToolArgs != nil {
		msg.ToolArgs = &model.JsonValue{Data: record.ToolArgs}
	}

	// 事务写入消息并更新会话统计
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		// 更新会话统计
		return tx.Model(&model.AgentSession{}).
			Where("id = ?", sessionId).
			Updates(map[string]any{
				"message_count": gorm.Expr("message_count + ?", 1),
				"total_tokens":  gorm.Expr("total_tokens + ?", record.Tokens),
			}).Error
	})
}

// findMessages 查询会话消息(公用底层查询)
func (s *AgentService) findMessages(sessionId int64) ([]model.AgentMessage, error) {
	var (
		m        model.AgentMessage
		messages []model.AgentMessage
		db       = s.DB(&m)
	)

	err := db.Model(&m).Where("session_id = ?", sessionId).Order("id ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// LoadHistory 加载会话历史(转换为agent消息供模型使用)
func (s *AgentService) LoadHistory(sessionId int64) ([]agent.Message, error) {
	messages, err := s.findMessages(sessionId)
	if err != nil {
		return nil, err
	}

	history := make([]agent.Message, 0, len(messages))
	for _, msg := range messages {
		history = append(history, s.toAgentMessage(msg))
	}
	return history, nil
}

// History 获取会话历史记录
func (s *AgentService) History(sessionId int64) ([]model.AgentMessage, error) {
	return s.findMessages(sessionId)
}

// toAgentMessage 转换数据库消息为agent消息
func (s *AgentService) toAgentMessage(msg model.AgentMessage) agent.Message {
	m := agent.Message{
		Role:       msg.Role,
		Content:    msg.Content,
		ToolCallId: msg.ToolCallId,
	}

	// assistant工具调用消息,重建tool_calls供模型使用
	if msg.Role == "assistant" && msg.ToolCallId != "" {
		m.ToolCalls = []*agent.MessageToolCall{
			{
				Id:   msg.ToolCallId,
				Type: "function",
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{
					Name:      msg.ToolName,
					Arguments: s.toolArgsToJson(msg.ToolArgs),
				},
			},
		}
	}

	return m
}

// toolArgsToJson 工具参数转JSON字符串
func (s *AgentService) toolArgsToJson(args *model.JsonValue) string {
	if args == nil || args.Data == nil {
		return "{}"
	}
	data, err := json.Marshal(args)
	if err != nil {
		return "{}"
	}
	return string(data)
}
