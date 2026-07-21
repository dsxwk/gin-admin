package model

import (
	"database/sql/driver"
	"fmt"
	"gin/pkg"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sort"
	"strings"
	"time"
)

const (
	CreatedField = "created_at"
	UpdatedField = "updated_at"
	DeletedField = "deleted_at"
)

type DateTime time.Time

func (t *DateTime) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte(`""`), nil
	}
	formatted := fmt.Sprintf("\"%s\"", time.Time(*t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t DateTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *DateTime) Scan(value interface{}) error {
	if value == nil {
		*t = DateTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = DateTime(v)
		return nil
	case []byte:
		tt, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = DateTime(tt)
		return nil
	case string:
		tt, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*t = DateTime(tt)
		return nil
	default:
		return fmt.Errorf("cannot convert %v to timestamp", value)
	}
}

type DeletedAt struct {
	gorm.DeletedAt
}

func (d DeletedAt) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02 15:04:05"))), nil
}

type JsonValue struct {
	Data any
}

// Scan 读取json
func (j *JsonValue) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("cannot scan %T into JSONValue", value)
	}

	if len(bytes) == 0 {
		j.Data = nil
		return nil
	}

	return json.Unmarshal(bytes, &j.Data)
}

// Value 写入json
func (j JsonValue) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// MarshalJSON 输出Data内容
func (j JsonValue) MarshalJSON() ([]byte, error) {
	// 不输出{ "Data": ... },直接输出内容
	return json.Marshal(j.Data)
}

type ArrayString []string

func (j ArrayString) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *ArrayString) Scan(value interface{}) error {
	if value == nil {
		*j = ArrayString{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("cannot scan %T into ArrayString", value)
	}
	return json.Unmarshal(bytes, j)
}

type ArrayInt64 []int64

func (j ArrayInt64) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *ArrayInt64) Scan(value interface{}) error {
	if value == nil {
		*j = ArrayInt64{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("cannot scan %T into ArrayInt64", value)
	}
	return json.Unmarshal(bytes, j)
}

// BatchUpdateSql 生成批量更新SQL语句
// model: 模型,例如&model.SystemConfig{}
// data: 更新数据,每个元素是一个包含列名和值的映射
// primaryKey: 主键字段名
// conditions: 附加的WHERE条件(可选)
// 返回值是生成的SQL语句和对应的参数值
// 调用示例: sql, values := BatchUpdateSql("table", data, "id", map[string]interface{}{"status": 1})
// facade.DB().Exec(sql, values...)
func BatchUpdateSql(db *gorm.DB, model any, data []map[string]interface{}, primaryKey string, conditions map[string]interface{}) (string, []interface{}) {
	if len(data) == 0 {
		return "", nil
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", nil
	}

	// 表名
	table := stmt.Schema.Table

	// 主键
	if primaryKey == "" {
		if stmt.Schema.PrioritizedPrimaryField == nil {
			return "", nil
		}
		primaryKey = stmt.Schema.PrioritizedPrimaryField.DBName
	} else {
		primaryKey = lo.SnakeCase(primaryKey)
	}

	var (
		values []interface{}
		ids    []interface{}
	)

	// db字段->原始key
	columnMap := make(map[string]string)

	for _, row := range data {
		id, ok := row[primaryKey]
		if !ok {
			// 支持id->ID
			id, ok = row["id"]
		}
		if !ok {
			id, ok = row["ID"]
		}
		if !ok {
			continue
		}

		ids = append(ids, id)
		for key := range row {
			dbName, _ok := getDBName(stmt.Schema, key)
			if !_ok {
				continue
			}

			if dbName == primaryKey {
				continue
			}

			columnMap[dbName] = key
		}
	}

	if len(ids) == 0 {
		return "", nil
	}

	columns := make([]string, 0, len(columnMap))
	for k := range columnMap {
		columns = append(columns, k)
	}
	sort.Strings(columns)

	var setSQL []string
	for _, dbColumn := range columns {
		originKey := columnMap[dbColumn]
		var builder strings.Builder
		builder.WriteString(
			fmt.Sprintf("`%s` = CASE `%s` ", dbColumn, primaryKey),
		)

		for _, row := range data {
			id := firstValue(row, primaryKey, "id", "ID")
			builder.WriteString("WHEN ? THEN ? ")
			values = append(values, id)

			if v, ok := row[originKey]; ok {
				values = append(values, v)
			} else {
				values = append(values, nil)
			}
		}

		builder.WriteString("END")
		setSQL = append(setSQL, builder.String())
	}

	// where id in
	in := make([]string, len(ids))
	for i := range ids {
		in[i] = "?"
		values = append(values, ids[i])
	}

	where := []string{
		fmt.Sprintf("`%s` IN (%s)", primaryKey, strings.Join(in, ",")),
	}

	for k, v := range conditions {
		dbName, ok := getDBName(stmt.Schema, k)
		if !ok {
			dbName = lo.SnakeCase(k)
		}

		where = append(where, fmt.Sprintf("`%s` = ?", dbName))
		values = append(values, v)
	}

	// 软删除
	if stmt.Schema.LookUpField("DeletedAt") != nil {
		where = append(where, "`"+DeletedField+"` IS NULL")
	}

	sql := pkg.Sprintf(
		"UPDATE `%s` SET %s WHERE %s",
		table,
		strings.Join(setSQL, ", "),
		strings.Join(where, " AND "),
	)

	return sql, values
}

func firstValue(m map[string]interface{}, keys ...string) interface{} {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			return v
		}
	}
	return nil
}

// 获取数据库字段名
func getDBName(s *schema.Schema, key string) (string, bool) {
	// Go字段
	if field, ok := s.FieldsByName[key]; ok {
		if _, relation := s.Relationships.Relations[field.Name]; relation {
			return "", false
		}
		return field.DBName, true
	}

	// 数据库字段
	if field, ok := s.FieldsByDBName[key]; ok {
		if _, relation := s.Relationships.Relations[field.Name]; relation {
			return "", false
		}
		return field.DBName, true
	}

	// snake_case
	snake := lo.SnakeCase(key)
	if field, ok := s.FieldsByDBName[snake]; ok {
		if _, relation := s.Relationships.Relations[field.Name]; relation {
			return "", false
		}
		return field.DBName, true
	}

	// json tag
	for _, field := range s.Fields {
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == key {
			if _, relation := s.Relationships.Relations[field.Name]; relation {
				return "", false
			}

			return field.DBName, true
		}
	}

	return "", false
}

// FilterFields 过滤非模型字段
func FilterFields(db *gorm.DB, model any, raw map[string]interface{}) map[string]interface{} {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return raw
	}

	filtered := make(map[string]interface{})

	for k, v := range raw {
		snakeKey := lo.SnakeCase(k)
		if field, ok := stmt.Schema.FieldsByDBName[snakeKey]; ok {
			// 检查是否为关联字段,通过字段名在Relationships中查找
			if _, ok = stmt.Schema.Relationships.Relations[field.Name]; ok {
				continue
			}
			filtered[snakeKey] = v
		}
	}

	return filtered
}
