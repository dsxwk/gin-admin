package cache

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

// encodeValue 编码缓存值为JSON
func encodeValue(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("缓存值JSON编码失败: %w", err)
	}

	return data, nil
}

// decodeValue 解析缓存JSON值
func decodeValue(value []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(value))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("缓存值JSON解析失败: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("缓存值包含多余内容")
	}

	return normalizeValue(result), nil
}

// normalizeValue 标准化JSON数字类型
func normalizeValue(value any) any {
	switch item := value.(type) {
	case json.Number:
		text := item.String()
		if strings.ContainsAny(text, ".eE") {
			number, err := item.Float64()
			if err == nil {
				return number
			}

			return text
		}

		number, err := strconv.ParseInt(text, 10, 64)
		if err == nil {
			return number
		}

		unsigned, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			return unsigned
		}

		return text
	case map[string]any:
		for key, value := range item {
			item[key] = normalizeValue(value)
		}

		return item
	case []any:
		for index, value := range item {
			item[index] = normalizeValue(value)
		}

		return item
	default:
		return value
	}
}
