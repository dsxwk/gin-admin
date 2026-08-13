package model

const TableNameAgentSession = "agent_session"

// AgentSession AI会话主表
type AgentSession struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	UserId       int64      `gorm:"column:user_id;not null;type:int(10) unsigned;comment:用户ID" json:"userId" form:"userId"`
	Title        string     `gorm:"column:title;not null;type:varchar(200);comment:会话标题" json:"title" form:"title"`
	Provider     string     `gorm:"column:provider;not null;type:varchar(50);comment:模型提供商,如deepseek、openai、moonshot" json:"provider" form:"provider"`
	Model        string     `gorm:"column:model;not null;type:varchar(50);comment:模型名,如deepseek-v4-pro" json:"model" form:"model"`
	MessageCount int64      `gorm:"column:message_count;not null;default:0;type:int(10) unsigned;comment:该会话总消息条数(含user、assistant、tool)" json:"messageCount" form:"messageCount"`
	TotalTokens  int64      `gorm:"column:total_tokens;not null;default:0;type:int(10) unsigned;comment:该会话累计消耗的token总数" json:"totalTokens" form:"totalTokens"`
	TraceId      string     `gorm:"column:trace_id;type:varchar(50);comment:请求追踪ID" json:"traceId" form:"traceId"`
	CreatedAt    *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt    *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt    *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*AgentSession) TableName() string {
	return TableNameAgentSession
}

// Connection 数据库连接名称
func (m *AgentSession) Connection() string {
	return "mysql"
}
