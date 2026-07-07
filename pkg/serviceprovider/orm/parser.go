package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm/schema"
)

// parse 递归解析查询条件
func parse(s *schema.Schema, filters map[string]any) (string, []any, error) {
	if len(filters) == 0 {
		return "", nil, nil
	}

	var (
		sqls []string
		args []any
	)
	for key, value := range filters {
		switch {
		case IsLogic(key):
			sql, values, err := parseGroup(s, key, value)
			if err != nil {
				return "", nil, err
			}

			if sql != "" {
				sqls = append(sqls, sql)
				args = append(args, values...)
			}

		case IsExists(key):
			sql, values, err := parseExists(s, key, value)
			if err != nil {
				return "", nil, err
			}

			if sql != "" {
				sqls = append(sqls, sql)
				args = append(args, values...)
			}

		default:
			sql, values, err := parseField(s, key, value)
			if err != nil {
				return "", nil, err
			}

			if sql != "" {
				sqls = append(sqls, sql)
				args = append(args, values...)
			}
		}
	}

	return JoinSQL(sqls, "AND"), args, nil
}

// parseGroup 解析AND和OR
func parseGroup(s *schema.Schema, logic string, value any) (string, []any, error) {
	items, ok := value.([]any)
	if !ok {
		return "", nil, fmt.Errorf("%s必须为数组", logic)
	}

	var (
		sqls []string
		args []any
	)

	for _, item := range items {
		m, _ok := item.(map[string]any)
		if !_ok {
			continue
		}

		sql, values, err := parse(s, m)
		if err != nil {
			return "", nil, err
		}

		if sql == "" {
			continue
		}

		sqls = append(sqls, Wrap(sql))
		args = append(args, values...)
	}

	if strings.EqualFold(logic, "or") {
		return JoinSQL(sqls, "OR"), args, nil
	}

	return JoinSQL(sqls, "AND"), args, nil
}

// parseExists 解析EXISTS
func parseExists(s *schema.Schema, key string, value any) (string, []any, error) {
	conditions, ok := value.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("%s必须为对象", key)
	}

	var (
		sqls []string
		args []any
	)
	not := strings.EqualFold(key, "not exist")
	for relation, condition := range conditions {
		filter, _ok := condition.(map[string]any)
		if !_ok {
			return "", nil, fmt.Errorf("%s查询条件必须为对象", relation)
		}

		sql, values, err := BuildExists(s, relation, filter, not)

		if err != nil {
			return "", nil, err
		}

		if sql == "" {
			continue
		}

		sqls = append(sqls, sql)
		args = append(args, values...)
	}

	return JoinSQL(sqls, "AND"), args, nil
}

// parseField 解析字段
func parseField(s *schema.Schema, field string, value any) (string, []any, error) {
	name, operator := SplitOperator(field)
	if IsRelation(name) {
		return BuildRelation(s, name, operator, value)
	}

	f := FindField(s, name)
	if f == nil {
		return "", nil, fmt.Errorf("字段[%s]不存在", name)
	}

	return BuildOperator(s.Table+"."+f.DBName, operator, value)
}
