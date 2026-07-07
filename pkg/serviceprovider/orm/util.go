package orm

import (
	"github.com/samber/lo"
	"gorm.io/gorm/schema"
	"regexp"
	"strings"
)

var (
	operatorRegexp = regexp.MustCompile(`\[(.+)]$`)
)

// SplitOperator 拆分字段和操作符
// 例如:
// name -> name,=
// age[>] -> age,>
// id[in] -> id,in
func SplitOperator(field string) (string, string) {
	match := operatorRegexp.FindStringSubmatch(field)
	if len(match) != 2 {
		return field, "="
	}

	return strings.TrimSpace(field[:len(field)-len(match[0])]),
		strings.ToLower(strings.TrimSpace(match[1]))
}

// FindField 查找字段
func FindField(s *schema.Schema, name string) *schema.Field {
	if s == nil || name == "" {
		return nil
	}

	if field := s.LookUpField(name); field != nil {
		return field
	}

	if field := s.LookUpField(lo.PascalCase(name)); field != nil {
		return field
	}

	if field := s.LookUpField(lo.SnakeCase(name)); field != nil {
		return field
	}

	return nil
}

// FindRelation 查找关联
func FindRelation(s *schema.Schema, name string) *schema.Relationship {
	if s == nil || name == "" {
		return nil
	}

	if rel, ok := s.Relationships.Relations[name]; ok {
		return rel
	}

	if rel, ok := s.Relationships.Relations[lo.PascalCase(name)]; ok {
		return rel
	}

	return nil
}

// IsRelation 判断是否关联字段
func IsRelation(field string) bool {
	return strings.Contains(field, ".")
}

// IsJSON 判断是否JSON字段
func IsJSON(field string) bool {
	return strings.Contains(field, "->")
}

// IsLogic 判断是否逻辑关键字
func IsLogic(key string) bool {
	switch strings.ToLower(key) {
	case "and", "or":
		return true
	default:
		return false
	}
}

// IsExists 判断是否Exists关键字
func IsExists(key string) bool {
	switch strings.ToLower(key) {
	case "exist", "not exist":
		return true
	default:
		return false
	}
}

// Wrap 包裹SQL
func Wrap(sql string) string {
	if sql == "" {
		return ""
	}
	return "(" + sql + ")"
}

// JoinSQL 拼接SQL
func JoinSQL(sqls []string, logic string) string {
	switch len(sqls) {
	case 0:
		return ""
	case 1:
		return sqls[0]
	default:
		return Wrap(strings.Join(sqls, " "+logic+" "))
	}
}
