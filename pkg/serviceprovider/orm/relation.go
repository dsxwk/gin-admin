package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm/schema"
)

// RelationInfo 关联信息
type RelationInfo struct {
	Relation *schema.Relationship
	Schema   *schema.Schema
	Field    *schema.Field
}

// ParseRelation 解析关联路径
// 例如: RoleMenus.Menu.Name
func ParseRelation(root *schema.Schema, path string) (*RelationInfo, error) {
	items := strings.Split(path, ".")
	if len(items) < 2 {
		return nil, fmt.Errorf("关联字段[%s]格式错误", path)
	}

	current := root
	var relation *schema.Relationship

	for _, name := range items[:len(items)-1] {
		relation = FindRelation(current, name)
		if relation == nil {
			return nil, fmt.Errorf("关联[%s]不存在", name)
		}

		current = relation.FieldSchema
	}

	field := FindField(current, items[len(items)-1])
	if field == nil {
		return nil, fmt.Errorf("字段[%s]不存在", items[len(items)-1])
	}

	return &RelationInfo{
		Relation: relation,
		Schema:   current,
		Field:    field,
	}, nil
}

// JoinSQL 生成关联条件
func (r *RelationInfo) JoinSQL(parentTable string) string {
	sqls := make([]string, 0, len(r.Relation.References))

	for _, ref := range r.Relation.References {
		sqls = append(sqls,
			fmt.Sprintf(
				"%s.%s=%s.%s",
				r.Schema.Table,
				ref.ForeignKey.DBName,
				parentTable,
				ref.PrimaryKey.DBName,
			),
		)
	}

	return strings.Join(sqls, " AND ")
}

// BuildRelation 构建关联字段查询
func BuildRelation(root *schema.Schema, field string, operator string, value any) (string, []any, error) {
	info, err := ParseRelation(root, field)
	if err != nil {
		return "", nil, err
	}

	expr, args, err := BuildOperator(info.Schema.Table+"."+info.Field.DBName, operator, value)
	if err != nil {
		return "", nil, err
	}

	sql := fmt.Sprintf("EXISTS(SELECT 1 FROM %s WHERE %s AND %s)", info.Schema.Table, info.JoinSQL(root.Table), expr)

	return sql, args, nil
}

// BuildExists 构建Exists查询
func BuildExists(root *schema.Schema, relation string, filter map[string]any, not bool) (string, []any, error) {
	info, err := ParseRelation(root, relation+".id")
	if err != nil {
		return "", nil, err
	}

	sql, args, err := parse(info.Schema, filter)
	if err != nil {
		return "", nil, err
	}

	where := info.JoinSQL(root.Table)

	if sql != "" {
		where += " AND " + sql
	}

	result := fmt.Sprintf("EXISTS(SELECT 1 FROM %s WHERE %s)", info.Schema.Table, where)

	if not {
		result = "NOT " + result
	}

	return result, args, nil
}
