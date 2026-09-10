package model

const TableNameOperatorLog = "operator_log"

// OperatorLog 操作日志表
type OperatorLog struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Ip         string     `gorm:"column:ip;not null;type:varchar(50);comment:ip地址" json:"ip" form:"ip"`
	Method     string     `gorm:"column:method;not null;type:varchar(20);comment:请求方式" json:"method" form:"method"`
	Header     string     `gorm:"column:header;not null;type:text;comment:请求头" json:"header" form:"header"`
	Uri        string     `gorm:"column:uri;not null;type:varchar(500);comment:路由地址" json:"uri" form:"uri"`
	Lang       string     `gorm:"column:lang;not null;type:varchar(20);comment:语言" json:"lang" form:"lang"`
	Params     *JsonValue `gorm:"column:params;type:json;comment:请求参数" json:"params" form:"params"`
	UserId     int64      `gorm:"column:user_id;not null;default:0;type:int(10) unsigned;comment:操作用户ID" json:"userId" form:"userId"`
	TraceId    string     `gorm:"column:trace_id;not null;type:varchar(50);comment:追踪ID" json:"traceId" form:"traceId"`
	StatusCode int64      `gorm:"column:status_code;not null;default:0;type:int(10);comment:响应状态码" json:"statusCode" form:"statusCode"`
	UserAgent  string     `gorm:"column:user_agent;not null;type:varchar(255);comment:用户代理" json:"userAgent" form:"userAgent"`
	CostMs     float64    `gorm:"column:cost_ms;not null;default:0.0000;type:float(10,4);comment:耗时(毫秒)" json:"costMs" form:"costMs"`
	User       *User      `gorm:"foreignKey:user_id;references:id" json:"user"`
	CreatedAt  *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt  *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt  *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*OperatorLog) TableName() string {
	return TableNameOperatorLog
}

// Connection 数据库连接名称
func (m *OperatorLog) Connection() string {
	return "mysql"
}
