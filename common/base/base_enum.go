package base

// Item 枚举项
type Item[T comparable] struct {
	Value T      `json:"value"`
	Desc  string `json:"desc"`
}

// Enum 通用枚举基类
type Enum[T comparable] struct {
	items []Item[T]
}

// NewEnum 创建枚举实例
func NewEnum[T comparable](items ...Item[T]) *Enum[T] {
	return &Enum[T]{items: items}
}

// Get 获取枚举列表
func (e *Enum[T]) Get() []Item[T] {
	return e.items
}

// Value 根据描述获取值
func (e *Enum[T]) Value(desc string) T {
	var zero T
	for _, item := range e.items {
		if item.Desc == desc {
			return item.Value
		}
	}
	return zero
}

// Desc 根据值获取描述
func (e *Enum[T]) Desc(value T) string {
	for _, item := range e.items {
		if item.Value == value {
			return item.Desc
		}
	}
	return ""
}

// ContainsValue 检查值是否存在
func (e *Enum[T]) ContainsValue(value T) bool {
	for _, item := range e.items {
		if item.Value == value {
			return true
		}
	}
	return false
}

// ContainsDesc 检查描述是否存在
func (e *Enum[T]) ContainsDesc(desc string) bool {
	for _, item := range e.items {
		if item.Desc == desc {
			return true
		}
	}
	return false
}

// Map 转换为map映射
func (e *Enum[T]) Map() map[T]string {
	m := make(map[T]string, len(e.items))
	for _, item := range e.items {
		m[item.Value] = item.Desc
	}
	return m
}

// Len 返回枚举数量
func (e *Enum[T]) Len() int {
	return len(e.items)
}
