package pkg

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// Sprintf 格式化字符串
func Sprintf(format string, args ...any) string {
	var b strings.Builder

	// 预估容量(可选)
	b.Grow(len(format) + 64)

	_, _ = fmt.Fprintf(&b, format, args...)

	return b.String()
}
