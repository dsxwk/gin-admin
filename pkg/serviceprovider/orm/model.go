package orm

// Model 数据库模型接口
type Model interface {
	TableName() string
	Connection() string
}
