package stringx

import (
	"fmt"
	"math/rand"
	"unicode"
)

// UcFirst 首字母大写
func UcFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LcFirst 首字母小写
func LcFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// StringLength 获取字符串长度
func StringLength(s string) int {
	return len([]rune(s))
}

// Spaces 返回指定数量的空格
func Spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

// RandString 生成指定长度的随机字符串
func RandString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
