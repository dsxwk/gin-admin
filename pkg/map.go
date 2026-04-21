package pkg

import (
	"fmt"
	"strings"
)

// HasKey 检查map键名是否存在(支持任意键类型,但是键类型必须是可比较类型)
func HasKey[K comparable, V any](data map[K]V, key K) bool {
	_, exists := data[key]
	return exists
}

// InArray 检查数组中是否存在某个值
func InArray[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// InArrayFast 检查数组中是否存在某个值(高性能适合多次查找场景)
func InArrayFast[T comparable](val T, arr []T) bool {
	m := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	_, ok := m[val]
	return ok
}

// ArrayFilter 过滤切片,保留满足条件的元素
func ArrayFilter[T any](arr []T, fn func(T) bool) []T {
	res := make([]T, 0)
	for _, v := range arr {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}

// FilterFields 过滤字段
// map[string]any, map[string]any → map[string]any
// []map[string]any, []map[string]any → []map[string]any
func FilterFields(data any, filter any) any {
	switch d := data.(type) {

	// map
	case map[string]any:
		f, ok := filter.(map[string]any)
		if !ok {
			return d // filter不是map,直接返回原数据
		}
		return filterMap(d, f)

	// slice
	case []map[string]any:
		f, ok := filter.([]map[string]any)
		if !ok {
			return d // filter不是slice,返回原数据
		}

		// slice长度不一致无法过滤
		if len(d) != len(f) {
			return d
		}

		result := make([]map[string]any, 0, len(d))
		for i := range d {
			result = append(result, filterMap(d[i], f[i]))
		}
		return result
	}

	return data
}

// filterMap 过滤map
func filterMap(data map[string]any, filter map[string]any) map[string]any {
	filtered := make(map[string]any)
	for k, v := range filter {
		if HasKey(data, k) {
			filtered[k] = v
		}
	}
	return filtered
}

// ArrayToString 将数组格式化为字符串
func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}
