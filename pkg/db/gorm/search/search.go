package search

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"unicode"
)

/**
Example:
`
/api/v1/user?page=1&pageSize=100&__search={
  "or": [
    {
      "and": [
        { "createdAt": [">", "2025-01-01"] },
        { "createdAt": ["<", "2026-01-01"] },
        { "not exist": { "UserRoles.name": "" } }
      ]
    },
    {
      "username": "admin"
    }
  ]
}
`
`
SELECT * FROM `user` WHERE ((((user.created_at > '2025-01-01') AND (user.created_at < '2026-01-01') AND (NOT EXISTS (SELECT 1 FROM user_roles WHERE user_roles.user_id = user.id AND user_roles.name = ''))) OR (user.username = 'admin'))) AND `user`.`deleted_at` IS NULL ORDER BY id DESC LIMIT 10
`
`
/api/v1/menu?page=1&pageSize=10&__search={
    "or": [
        {
            "and": [
                {
                    "createdAt": [
                        ">",
                        "2025-01-01"
                    ]
                },
                {
                    "createdAt": [
                        "<",
                        "2026-01-01"
                    ]
                },
                {
                    "name": ""
                },
                {
                    "$.meta.icon": [
                        "=",
                        "ele-Collection"
                    ]
                }
            ]
        }
    ]
}
`
`
db := DB.Model(&model.Menu{})
if _search != nil {
	whereSql, args, _ := search.BuildCondition(_search, db, model.Menu{})

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}
}
`
`
SELECT * FROM `menu` WHERE ((((menu.created_at > '2025-01-01') AND (menu.created_at < '2026-01-01') AND (menu.name = '') AND (JSON_EXTRACT(meta, '$.icon') = 'ele-Collection')))) AND `menu`.`deleted_at` IS NULL ORDER BY id DESC LIMIT 10
`
*/

// BuildCondition 构建动态sql条件
// filters: map[string]interface{}格式参考示例
// model: 必须是非切片类型模型,用于解析schema
func BuildCondition(filters map[string]interface{}, db *gorm.DB, model interface{}) (string, []interface{}, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", nil, err
	}

	sql, args, err := parseLogic(filters, stmt.Schema, db)
	if err != nil {
		return "", nil, err
	}
	if sql != "" {
		sql = "(" + sql + ")"
	}
	return sql, args, nil
}

// parseLogic 递归处理解析and/or/exist/not exist逻辑
func parseLogic(filters map[string]interface{}, s *schema.Schema, db *gorm.DB) (string, []interface{}, error) {
	var parts []string
	var args []interface{}

	for k, v := range filters {
		key := strings.ToLower(k)
		switch key {
		case "and", "or":
			items, ok := v.([]interface{})
			if !ok {
				continue
			}
			var subParts []string
			var subArgs []interface{}
			for _, it := range items {
				if m, _ok := it.(map[string]interface{}); _ok {
					sql, a, err := parseLogic(m, s, db)
					if err != nil {
						return "", nil, err
					}
					if sql != "" {
						subParts = append(subParts, "("+sql+")")
						subArgs = append(subArgs, a...)
					}
				}
			}
			if len(subParts) > 0 {
				parts = append(parts, strings.Join(subParts, " "+strings.ToUpper(key)+" "))
				args = append(args, subArgs...)
			}
		case "exist", "not exist":
			m, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			for relPath, cond := range m {
				sql, a, err := buildExistCondition(s, db, relPath, cond, key == "exist")
				if err != nil {
					return "", nil, err
				}
				if sql != "" {
					parts = append(parts, sql)
					args = append(args, a...)
				}
			}
		default:
			sql, a := buildFieldSQL(db, s, k, v)
			if sql != "" {
				parts = append(parts, sql)
				args = append(args, a...)
			}
		}
	}

	return strings.Join(parts, " AND "), args, nil
}

// buildExistCondition 构建exists/not exists子查询
func buildExistCondition(parent *schema.Schema, db *gorm.DB, path string, cond interface{}, positive bool) (string, []interface{}, error) {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("exist must be relation.field")
	}

	curSchema := parent
	var lastRel *schema.Relationship

	for i := 0; i < len(parts)-1; i++ {
		name := toPascalCase(parts[i])
		rel, ok := curSchema.Relationships.Relations[name]
		if !ok {
			return "", nil, fmt.Errorf("relation %s not found on %s", name, curSchema.Name)
		}
		lastRel = rel
		curSchema = rel.FieldSchema
	}

	subTable := curSchema.Table
	fieldName := parts[len(parts)-1]
	f := curSchema.LookUpField(fieldName)
	if f == nil {
		f = curSchema.LookUpField(toPascalCase(fieldName))
	}
	if f == nil {
		return "", nil, nil
	}

	if lastRel == nil {
		return "", nil, nil
	}

	// 构建join条件
	var joins []string
	for _, ref := range lastRel.References {
		joins = append(joins, fmt.Sprintf("%s.%s = %s.%s", subTable, ref.ForeignKey.DBName, parent.Table, ref.PrimaryKey.DBName))
	}

	where := strings.Join(joins, " AND ")
	args := []interface{}{}

	// 条件支持map[string]interface{}或简单值
	if m, ok := cond.(map[string]interface{}); ok {
		for k, v := range m {
			if sub, _ok := v.(map[string]interface{}); _ok {
				sql, a, err := buildExistCondition(curSchema, db, k, sub, true)
				if err != nil {
					return "", nil, err
				}
				where += " AND " + sql
				args = append(args, a...)
			} else {
				f2 := curSchema.LookUpField(k)
				if f2 == nil {
					return "", nil, fmt.Errorf("field %s not found on %s", k, curSchema.Name)
				}
				where += fmt.Sprintf(" AND %s.%s = ?", subTable, f2.DBName)
				args = append(args, v)
			}
		}
	} else {
		where += fmt.Sprintf(" AND %s.%s = ?", subTable, f.DBName)
		args = append(args, cond)
	}

	sql := fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s)", subTable, where)
	if !positive {
		sql = "NOT " + sql
	}

	return sql, args, nil
}

// buildFieldSQL 构建普通字段sql,支持json、关联字段、各种操作符
func buildFieldSQL(db *gorm.DB, s *schema.Schema, field string, value interface{}) (string, []interface{}) {
	op := "="
	val := value

	if arr, ok := value.([]interface{}); ok && len(arr) >= 2 {
		op = fmt.Sprint(arr[0])
		val = arr[1]
	}

	// json字段
	if strings.HasPrefix(field, "$.") || strings.HasPrefix(field, "json->") {
		return buildJsonSql(field, op, val)
	}

	// 关联字段a.b→exist子查询
	if strings.Contains(field, ".") {
		return buildRelationFieldSql(db, field, op, val)
	}

	// 普通字段
	f := s.LookUpField(field)
	if f == nil {
		f = s.LookUpField(toPascalCase(field))
	}
	if f == nil {
		// 字段不存在,直接忽略
		return "", nil
	}
	col := s.Table + "." + f.DBName
	op = strings.ToLower(op)

	switch op {
	case "=", "eq":
		return fmt.Sprintf("%s = ?", col), []interface{}{val}
	case "!=", "<>", "ne":
		return fmt.Sprintf("%s <> ?", col), []interface{}{val}
	case ">", "gt":
		return fmt.Sprintf("%s > ?", col), []interface{}{val}
	case ">=", "gte":
		return fmt.Sprintf("%s >= ?", col), []interface{}{val}
	case "<", "lt":
		return fmt.Sprintf("%s < ?", col), []interface{}{val}
	case "<=", "lte":
		return fmt.Sprintf("%s <= ?", col), []interface{}{val}
	case "in":
		return fmt.Sprintf("%s IN (?)", col), []interface{}{val}
	case "not in":
		return fmt.Sprintf("%s NOT IN (?)", col), []interface{}{val}
	case "is null":
		return fmt.Sprintf("%s IS NULL", col), nil
	case "is not null":
		return fmt.Sprintf("%s IS NOT NULL", col), nil
	case "like":
		return fmt.Sprintf("%s LIKE ?", col), []interface{}{"%" + fmt.Sprint(val) + "%"}
	case "left like":
		return fmt.Sprintf("%s LIKE ?", col), []interface{}{fmt.Sprint(val) + "%"}
	case "right like":
		return fmt.Sprintf("%s LIKE ?", col), []interface{}{"%" + fmt.Sprint(val)}
	case "between":
		arr, ok := val.([]interface{})
		if ok && len(arr) == 2 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", col), []interface{}{arr[0], arr[1]}
		}
	case "not between":
		arr, ok := val.([]interface{})
		if ok && len(arr) == 2 {
			return fmt.Sprintf("%s NOT BETWEEN ? AND ?", col), []interface{}{arr[0], arr[1]}
		}
	}

	return fmt.Sprintf("%s = ?", col), []interface{}{val}
}

// buildRelationFieldSql 将关联字段a.b转化为exists子查询
func buildRelationFieldSql(db *gorm.DB, field string, op string, val interface{}) (string, []interface{}) {
	parts := strings.Split(field, ".")
	if len(parts) < 2 {
		return "", nil
	}

	rootSchema := db.Statement.Schema
	curSchema := rootSchema
	var rel *schema.Relationship

	for i := 0; i < len(parts)-1; i++ {
		name := toPascalCase(parts[i])
		r, ok := curSchema.Relationships.Relations[name]
		if !ok {
			return "", nil
		}
		rel = r
		curSchema = r.FieldSchema
	}

	if rel == nil {
		return "", nil
	}

	subTable := curSchema.Table
	subColumn := camelToSnake(parts[len(parts)-1])

	var joins []string
	for _, ref := range rel.References {
		left := fmt.Sprintf("%s.%s", subTable, ref.ForeignKey.DBName)
		right := fmt.Sprintf("%s.%s", rootSchema.Table, ref.PrimaryKey.DBName)
		joins = append(joins, left+" = "+right)
	}

	sub := fmt.Sprintf("SELECT 1 FROM %s WHERE %s AND %s %s ?", subTable, strings.Join(joins, " AND "), subColumn, op)

	if op != "=" && (op == "!=" || op == "<>") {
		return "NOT EXISTS (" + sub + ")", []interface{}{val}
	}
	return "EXISTS (" + sub + ")", []interface{}{val}
}

// buildJsonSql json查询构建
func buildJsonSql(field, op string, val interface{}) (string, []interface{}) {
	var col, path string

	if strings.HasPrefix(field, "$.") {
		parts := strings.Split(field, ".")
		col = camelToSnake(parts[1])
		path = "$." + strings.Join(parts[2:], ".")
	}
	if strings.HasPrefix(field, "json->") {
		f := strings.TrimPrefix(field, "json->")
		parts := strings.Split(f, ".")
		col = camelToSnake(parts[0])
		path = "$." + strings.Join(parts[1:], ".")
	}

	op = strings.ToLower(op)
	switch op {
	case "json_contains":
		return fmt.Sprintf("JSON_CONTAINS(%s, ?)", col), []interface{}{val}
	case "in":
		return fmt.Sprintf("JSON_EXTRACT(%s, '%s') IN (?)", col, path), []interface{}{val}
	default:
		return fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ?", col, path), []interface{}{val}
	}
}

// toPascalCase 转大驼峰命名
func toPascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "_") {
		s = snakeToCamel(s)
	}
	if unicode.IsLower(rune(s[0])) {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// 将驼峰命名转为下划线命名
func camelToSnake(s string) string {
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

// 将下划线命名转为驼峰命名
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}
	return strings.Join(parts, "")
}
