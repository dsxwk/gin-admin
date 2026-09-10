package es

// SearchResponse 搜索响应
type SearchResponse[T any] struct {
	Took     int           `json:"took"`      // 耗时
	TimedOut bool          `json:"timed_out"` // 是否超时
	Hits     SearchHits[T] `json:"hits"`      // 命中结果
}

// SearchHits 搜索命中
type SearchHits[T any] struct {
	Total SearchTotal    `json:"total"` // 总数
	Hits  []SearchHit[T] `json:"hits"`  // 命中列表
}

// SearchTotal 搜索总数
type SearchTotal struct {
	Value    int64  `json:"value"`    // 总数
	Relation string `json:"relation"` // 关系
}

// SearchHit 单条命中
type SearchHit[T any] struct {
	Index  string  `json:"_index"`  // 索引
	ID     string  `json:"_id"`     // ID
	Score  float64 `json:"_score"`  // 得分
	Source T       `json:"_source"` // 文档数据
}

// DocumentResponse 文档响应
type DocumentResponse[T any] struct {
	Index  string `json:"_index"`  // 索引
	ID     string `json:"_id"`     // ID
	Found  bool   `json:"found"`   // 是否存在
	Source T      `json:"_source"` // 文档数据
}
