package pkg

import "strconv"

// Integer 整数类型
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// StringToInt 字符串转整数
func StringToInt[T Integer](s string) (T, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		var zero T
		return zero, err
	}
	return T(val), nil
}

// IntToString 整数转字符串
func IntToString[T Integer](v T) string {
	return strconv.FormatInt(int64(v), 10)
}
