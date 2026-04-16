package pkg

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// SnakeToCamel 将下划线命名转换为驼峰命名
func SnakeToCamel(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	upperNext := true
	for _, r := range s {
		if r == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			b.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// CamelToSnake 将驼峰命名转换为下划线命名
func CamelToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ToLowerCamel 将下划线命名转换为小驼峰命名
func ToLowerCamel(s string) string {
	if s == "" {
		return ""
	}
	camel := SnakeToCamel(s)
	return strings.ToLower(camel[:1]) + camel[1:]
}

// ToUpperCamel 将下划线命名转换为大驼峰命名
func ToUpperCamel(s string) string {
	c := cases.Title(language.Und)
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = c.String(parts[i])
	}
	return strings.Join(parts, "")
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
