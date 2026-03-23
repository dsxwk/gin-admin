package pkg

import (
	"fmt"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
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

// ArrayUnique 数组去重
func ArrayUnique[T comparable](arr []T) []T {
	m := make(map[T]struct{})
	res := make([]T, 0, len(arr))
	for _, v := range arr {
		if _, exists := m[v]; !exists {
			m[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}

// IndexOf 返回元素在切片中的下标,未找到返回-1
func IndexOf[T comparable](val T, arr []T) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
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

// StructToMap 将结构体转换为map
func StructToMap[T any](v T) any {
	var (
		m map[string]any
	)

	b, _ := json.Marshal(v)
	_ = json.Unmarshal(b, &m)
	return m
}

// ArrayToString 将数组格式化为字符串
func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

// FilterModelFields 过滤非模型字段
func FilterModelFields(db *gorm.DB, model any, raw map[string]interface{}) map[string]interface{} {
	stmt := &gorm.Statement{DB: db}
	_ = stmt.Parse(model)

	filtered := make(map[string]interface{})

	for k, v := range raw {
		if _, ok := stmt.Schema.FieldsByDBName[CamelToSnake(k)]; ok {
			filtered[CamelToSnake(k)] = v
		}
	}

	return filtered
}
