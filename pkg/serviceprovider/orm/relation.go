package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm/schema"
)

// JoinStep 关联链中的一步
type JoinStep struct {
	Relation *schema.Relationship
	Schema   *schema.Schema
}

// RelationInfo 关联信息
type RelationInfo struct {
	Chain  []JoinStep
	Schema *schema.Schema
	Field  *schema.Field
}

// ParseRelation 解析关联路径
// 例如: deptLeaders.leader.fullName
func ParseRelation(root *schema.Schema, path string) (*RelationInfo, error) {
	items := strings.Split(path, ".")
	if len(items) < 2 {
		return nil, fmt.Errorf("关联字段[%s]格式错误", path)
	}

	current := root
	var chain []JoinStep

	for _, name := range items[:len(items)-1] {
		relation := FindRelation(current, name)
		if relation == nil {
			return nil, fmt.Errorf("关联[%s]不存在", name)
		}

		chain = append(chain, JoinStep{
			Relation: relation,
			Schema:   current,
		})

		current = relation.FieldSchema
	}

	field := FindField(current, items[len(items)-1])
	if field == nil {
		return nil, fmt.Errorf("字段[%s]不存在", items[len(items)-1])
	}

	return &RelationInfo{
		Chain:  chain,
		Schema: current,
		Field:  field,
	}, nil
}

// fromClause 生成FROM子句(含JOIN)
func (r *RelationInfo) fromClause() string {
	if len(r.Chain) == 0 {
		return r.Schema.Table
	}

	var parts []string
	// 第一个表
	first := r.Chain[0]
	parts = append(parts, first.Relation.FieldSchema.Table)

	// 后续表通过JOIN连接
	for i := 1; i < len(r.Chain); i++ {
		curr := r.Chain[i].Relation.FieldSchema.Table
		rel := r.Chain[i].Relation

		var joins []string
		for _, ref := range rel.References {
			joins = append(joins,
				fmt.Sprintf("%s.%s = %s.%s",
					ref.PrimaryKey.Schema.Table,
					ref.PrimaryKey.DBName,
					ref.ForeignKey.Schema.Table,
					ref.ForeignKey.DBName,
				),
			)
		}
		parts = append(parts, fmt.Sprintf("INNER JOIN %s ON %s", curr, strings.Join(joins, " AND ")))
	}

	return strings.Join(parts, " ")
}

// joinConditions 生成根表与第一个关联表的连接条件
func (r *RelationInfo) joinConditions(rootTable string) string {
	if len(r.Chain) == 0 {
		return ""
	}

	var conditions []string
	first := r.Chain[0]
	for _, ref := range first.Relation.References {
		conditions = append(conditions,
			fmt.Sprintf("%s.%s = %s.%s",
				ref.ForeignKey.Schema.Table,
				ref.ForeignKey.DBName,
				rootTable,
				ref.PrimaryKey.DBName,
			),
		)
	}

	return strings.Join(conditions, " AND ")
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

	from := info.fromClause()
	conditions := info.joinConditions(root.Table)

	var where string
	if conditions != "" {
		where = conditions + " AND " + expr
	} else {
		where = expr
	}

	sql := fmt.Sprintf("EXISTS(SELECT 1 FROM %s WHERE %s)", from, where)

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

	conditions := info.joinConditions(root.Table)

	if sql != "" {
		if conditions != "" {
			conditions += " AND " + sql
		} else {
			conditions = sql
		}
	}

	from := info.fromClause()
	result := fmt.Sprintf("EXISTS(SELECT 1 FROM %s WHERE %s)", from, conditions)

	if not {
		result = "NOT " + result
	}

	return result, args, nil
}
