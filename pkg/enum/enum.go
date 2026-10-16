package enum

import "slices"

// Item 枚举项
type Item[T comparable] struct {
	Value T      `json:"value"`
	Desc  string `json:"desc"`
}

// Enum 通用枚举
type Enum[T comparable] struct {
	items []Item[T]
}

// NewEnum 创建枚举实例
func NewEnum[T comparable](items ...Item[T]) *Enum[T] {
	return &Enum[T]{items: items}
}

// Get 获取枚举列表
func (e *Enum[T]) Get() []Item[T] {
	return slices.Clone(e.items)
}

// Value 根据描述获取值
func (e *Enum[T]) Value(desc string) T {
	value, _ := e.ValueOK(desc)
	return value
}

// ValueOK 根据描述获取值和是否存在
func (e *Enum[T]) ValueOK(desc string) (T, bool) {
	var zero T
	for _, item := range e.items {
		if item.Desc == desc {
			return item.Value, true
		}
	}
	return zero, false
}

// Desc 根据值获取描述
func (e *Enum[T]) Desc(value T) string {
	desc, _ := e.DescOK(value)
	return desc
}

// DescOK 根据值获取描述和是否存在
func (e *Enum[T]) DescOK(value T) (string, bool) {
	for _, item := range e.items {
		if item.Value == value {
			return item.Desc, true
		}
	}
	return "", false
}

// ContainsValue 检查值是否存在
func (e *Enum[T]) ContainsValue(value T) bool {
	_, ok := e.DescOK(value)
	return ok
}

// ContainsDesc 检查描述是否存在
func (e *Enum[T]) ContainsDesc(desc string) bool {
	_, ok := e.ValueOK(desc)
	return ok
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
