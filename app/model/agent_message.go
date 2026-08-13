package model

const TableNameAgentMessage = "agent_message"

// AgentMessage AI消息详情表
type AgentMessage struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	SessionId  int64      `gorm:"column:session_id;not null;type:int(10) unsigned;comment:AI会话ID" json:"sessionId" form:"sessionId"`
	Role       string     `gorm:"column:role;not null;type:varchar(20);comment:user=用户问题 assistant=模型回答 tool=工具执行结果 system=系统提示语" json:"role" form:"role"`
	Content    string     `gorm:"column:content;type:text;comment:消息正文,user存问题原文,assistant存回答或空(有工具调用时content为空),tool存工具返回的JSON" json:"content" form:"content"`
	ToolName   string     `gorm:"column:tool_name;type:varchar(100);comment:工具名称,assistant消息存要调用的工具名(如route:list),tool消息存已执行的工具名" json:"toolName" form:"toolName"`
	ToolCallId string     `gorm:"column:tool_call_id;type:varchar(100);comment:DeepSeek返回的工具调用ID,assistant和对应的tool消息共用一个ID,用于关联\"调了什么→返回了什么\"" json:"toolCallId" form:"toolCallId"`
	ToolArgs   *JsonValue `gorm:"column:tool_args;type:json;comment:工具参数,仅assistant消息,存模型传过来的参数JSON(如{\"command\":\"route:list\"})" json:"toolArgs" form:"toolArgs"`
	Tokens     int64      `gorm:"column:tokens;not null;default:0;type:int(10) unsigned;comment:本条消息消耗的token数（API返回的usage）" json:"tokens" form:"tokens"`
	CostMs     float64    `gorm:"column:cost_ms;not null;default:0.00;type:decimal(10,2) unsigned;comment:本轮对话耗时(毫秒),包含API调用+工具执行的总时间" json:"costMs" form:"costMs"`
	CreatedAt  *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt  *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt  *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*AgentMessage) TableName() string {
	return TableNameAgentMessage
}

// Connection 数据库连接名称
func (m *AgentMessage) Connection() string {
	return "mysql"
}
