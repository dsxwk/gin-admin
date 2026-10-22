package migration

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type columnKind int

const (
	columnID columnKind = iota
	columnString
	columnText
	columnInt
	columnBigInt
	columnBool
	columnDateTime
)

type indexDefinition struct {
	name    string
	columns []string
	unique  bool
}

// Table 数据表构建器
type Table struct {
	db      *gorm.DB
	name    string
	columns []*Column
	indexes []indexDefinition
	err     error
}

// Column 数据列构建器
type Column struct {
	table         *Table
	name          string
	kind          columnKind
	length        int
	nullable      bool
	unsigned      bool
	primary       bool
	autoIncrement bool
	unique        bool
	hasDefault    bool
	defaultValue  any
	comment       string
}

// CreateTable 创建数据表
func CreateTable(db *gorm.DB, name string, define func(table *Table)) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	if define == nil {
		return fmt.Errorf("表结构定义不能为空")
	}

	table := &Table{db: db, name: strings.TrimSpace(name)}
	if table.name == "" {
		return fmt.Errorf("表名不能为空")
	}
	define(table)
	if table.err != nil {
		return table.err
	}
	if len(table.columns) == 0 {
		return fmt.Errorf("表%s没有定义字段", table.name)
	}

	statement, err := table.createStatement()
	if err != nil {
		return err
	}
	if err = db.Exec(statement).Error; err != nil {
		return fmt.Errorf("创建表%s失败: %w", table.name, err)
	}

	for _, index := range table.indexes {
		if err = table.createIndex(index); err != nil {
			return err
		}
	}
	return nil
}

// DropTable 删除数据表
func DropTable(db *gorm.DB, name string) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("表名不能为空")
	}
	if err := db.Migrator().DropTable(name); err != nil {
		return fmt.Errorf("删除表%s失败: %w", name, err)
	}
	return nil
}

// Id 添加自增主键
func (t *Table) Id() *Column {
	column := t.addColumn("id", columnID)
	if column == nil {
		return &Column{table: t}
	}
	column.primary = true
	column.autoIncrement = true
	column.unsigned = true
	return column
}

// String 添加字符串字段
func (t *Table) String(name string, length int) *Column {
	column := t.addColumn(name, columnString)
	if column == nil {
		return &Column{table: t}
	}
	if length <= 0 {
		length = 255
	}
	column.length = length
	return column
}

// Text 添加文本字段
func (t *Table) Text(name string) *Column {
	return t.addColumn(name, columnText)
}

// Int 添加整数字段
func (t *Table) Int(name string) *Column {
	return t.addColumn(name, columnInt)
}

// BigInt 添加长整数字段
func (t *Table) BigInt(name string) *Column {
	return t.addColumn(name, columnBigInt)
}

// Boolean 添加布尔字段
func (t *Table) Boolean(name string) *Column {
	return t.addColumn(name, columnBool)
}

// DateTime 添加时间字段
func (t *Table) DateTime(name string) *Column {
	return t.addColumn(name, columnDateTime)
}

// Timestamps 添加创建和更新时间
func (t *Table) Timestamps() {
	t.DateTime("created_at").Nullable()
	t.DateTime("updated_at").Nullable()
}

// SoftDeletes 添加软删除时间
func (t *Table) SoftDeletes() {
	t.DateTime("deleted_at").Nullable()
}

// Index 添加普通索引
func (t *Table) Index(name string, columns ...string) *Table {
	t.addIndex(name, false, columns...)
	return t
}

// UniqueIndex 添加唯一索引
func (t *Table) UniqueIndex(name string, columns ...string) *Table {
	t.addIndex(name, true, columns...)
	return t
}

// NotNull 设置字段非空
func (c *Column) NotNull() *Column {
	c.nullable = false
	return c
}

// Nullable 设置字段可空
func (c *Column) Nullable() *Column {
	c.nullable = true
	return c
}

// Unsigned 设置无符号数值
func (c *Column) Unsigned() *Column {
	c.unsigned = true
	return c
}

// Unique 设置字段唯一
func (c *Column) Unique() *Column {
	c.unique = true
	return c
}

// Default 设置字段默认值
func (c *Column) Default(value any) *Column {
	c.hasDefault = true
	c.defaultValue = value
	return c
}

// Comment 设置字段注释
func (c *Column) Comment(comment string) *Column {
	c.comment = strings.TrimSpace(comment)
	return c
}

// addColumn 添加字段
func (t *Table) addColumn(name string, kind columnKind) *Column {
	name = strings.TrimSpace(name)
	if name == "" {
		t.setError("字段名不能为空")
		return nil
	}
	for _, column := range t.columns {
		if column.name == name {
			t.setError("字段重复:" + name)
			return nil
		}
	}

	column := &Column{table: t, name: name, kind: kind}
	t.columns = append(t.columns, column)
	return column
}

// addIndex 添加索引
func (t *Table) addIndex(name string, unique bool, columns ...string) {
	name = strings.TrimSpace(name)
	if name == "" {
		t.setError("索引名不能为空")
		return
	}
	if len(columns) == 0 {
		columns = t.inferIndexColumns(name)
	}
	if len(columns) == 0 {
		t.setError("索引" + name + "没有字段")
		return
	}
	t.indexes = append(t.indexes, indexDefinition{
		name:    name,
		columns: columns,
		unique:  unique,
	})
}

// inferIndexColumns 推断索引字段
func (t *Table) inferIndexColumns(name string) []string {
	candidate := strings.TrimPrefix(strings.TrimPrefix(name, "idx_"), "uniq_")
	for _, column := range t.columns {
		if column.name == candidate {
			return []string{candidate}
		}
	}
	if len(t.columns) > 0 {
		return []string{t.columns[len(t.columns)-1].name}
	}
	return nil
}

// setError 保存构建错误
func (t *Table) setError(message string) {
	if t.err == nil {
		t.err = fmt.Errorf("%s", message)
	}
}

// createStatement 生成建表语句
func (t *Table) createStatement() (string, error) {
	definitions := make([]string, 0, len(t.columns))
	for _, column := range t.columns {
		definition, err := t.columnDefinition(column)
		if err != nil {
			return "", err
		}
		definitions = append(definitions, definition)
	}

	return "CREATE TABLE " + t.quote(t.name) +
		" (" + strings.Join(definitions, ", ") + ")", nil
}

// columnDefinition 生成字段定义
func (t *Table) columnDefinition(column *Column) (string, error) {
	dialect := t.db.Dialector.Name()
	parts := []string{t.quote(column.name), t.columnType(dialect, column)}
	if column.primary {
		parts = append(parts, "PRIMARY KEY")
	}
	if !column.nullable {
		parts = append(parts, "NOT NULL")
	}
	if column.unique {
		parts = append(parts, "UNIQUE")
	}
	if column.hasDefault {
		parts = append(parts, "DEFAULT "+formatDefault(dialect, column.defaultValue))
	}
	if column.comment != "" && dialect == "mysql" {
		parts = append(parts, "COMMENT "+quoteSQLString(column.comment))
	}
	return strings.Join(parts, " "), nil
}

// columnType 获取字段数据库类型
func (t *Table) columnType(dialect string, column *Column) string {
	switch dialect {
	case "mysql":
		return t.mysqlColumnType(column)
	case "postgres":
		return t.postgresColumnType(column)
	case "sqlserver":
		return t.sqlserverColumnType(column)
	case "oracle":
		return t.oracleColumnType(column)
	default:
		return t.sqliteColumnType(column)
	}
}

// mysqlColumnType 获取MySQL字段类型
func (t *Table) mysqlColumnType(column *Column) string {
	switch column.kind {
	case columnID:
		return "BIGINT UNSIGNED AUTO_INCREMENT"
	case columnString:
		return fmt.Sprintf("VARCHAR(%d)", column.length)
	case columnText:
		return "TEXT"
	case columnInt:
		if column.unsigned {
			return "INT UNSIGNED"
		}
		return "INT"
	case columnBigInt:
		if column.unsigned {
			return "BIGINT UNSIGNED"
		}
		return "BIGINT"
	case columnBool:
		return "TINYINT(1)"
	case columnDateTime:
		return "DATETIME"
	default:
		return "VARCHAR(255)"
	}
}

// postgresColumnType 获取PostgreSQL字段类型
func (t *Table) postgresColumnType(column *Column) string {
	switch column.kind {
	case columnID:
		return "BIGSERIAL"
	case columnString:
		return fmt.Sprintf("VARCHAR(%d)", column.length)
	case columnText:
		return "TEXT"
	case columnInt:
		return "INTEGER"
	case columnBigInt:
		return "BIGINT"
	case columnBool:
		return "BOOLEAN"
	case columnDateTime:
		return "TIMESTAMP"
	default:
		return "VARCHAR(255)"
	}
}

// sqlserverColumnType 获取SQLServer字段类型
func (t *Table) sqlserverColumnType(column *Column) string {
	switch column.kind {
	case columnID:
		return "BIGINT IDENTITY(1,1)"
	case columnString:
		return fmt.Sprintf("NVARCHAR(%d)", column.length)
	case columnText:
		return "NVARCHAR(MAX)"
	case columnInt:
		return "INT"
	case columnBigInt:
		return "BIGINT"
	case columnBool:
		return "BIT"
	case columnDateTime:
		return "DATETIME2"
	default:
		return "NVARCHAR(255)"
	}
}

// oracleColumnType 获取Oracle字段类型
func (t *Table) oracleColumnType(column *Column) string {
	switch column.kind {
	case columnID:
		return "NUMBER(19) GENERATED BY DEFAULT AS IDENTITY"
	case columnString:
		return fmt.Sprintf("VARCHAR2(%d)", column.length)
	case columnText:
		return "CLOB"
	case columnInt:
		return "NUMBER(10)"
	case columnBigInt:
		return "NUMBER(19)"
	case columnBool:
		return "NUMBER(1)"
	case columnDateTime:
		return "TIMESTAMP"
	default:
		return "VARCHAR2(255)"
	}
}

// sqliteColumnType 获取SQLite字段类型
func (t *Table) sqliteColumnType(column *Column) string {
	switch column.kind {
	case columnID:
		return "INTEGER"
	case columnString:
		return fmt.Sprintf("VARCHAR(%d)", column.length)
	case columnText:
		return "TEXT"
	case columnInt, columnBigInt, columnBool:
		return "INTEGER"
	case columnDateTime:
		return "DATETIME"
	default:
		return "VARCHAR(255)"
	}
}

// createIndex 创建索引
func (t *Table) createIndex(index indexDefinition) error {
	columns := make([]string, 0, len(index.columns))
	for _, column := range index.columns {
		columns = append(columns, t.quote(column))
	}

	prefix := "CREATE INDEX "
	if index.unique {
		prefix = "CREATE UNIQUE INDEX "
	}
	statement := prefix + t.quote(index.name) +
		" ON " + t.quote(t.name) +
		" (" + strings.Join(columns, ", ") + ")"
	if err := t.db.Exec(statement).Error; err != nil {
		return fmt.Errorf("创建索引%s失败: %w", index.name, err)
	}
	return nil
}

// quote 引用数据库标识符
func (t *Table) quote(name string) string {
	switch t.db.Dialector.Name() {
	case "mysql":
		return "`" + strings.ReplaceAll(name, "`", "``") + "`"
	case "sqlserver":
		return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
	default:
		return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
	}
}

// formatDefault 格式化默认值
func formatDefault(dialect string, value any) string {
	switch typed := value.(type) {
	case nil:
		return "NULL"
	case string:
		return quoteSQLString(typed)
	case bool:
		if dialect == "postgres" {
			if typed {
				return "TRUE"
			}
			return "FALSE"
		}
		if typed {
			return "1"
		}
		return "0"
	case time.Time:
		return quoteSQLString(typed.Format("2006-01-02 15:04:05"))
	case int:
		return strconv.Itoa(typed)
	case int8:
		return strconv.FormatInt(int64(typed), 10)
	case int16:
		return strconv.FormatInt(int64(typed), 10)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case int64:
		return strconv.FormatInt(typed, 10)
	case uint:
		return strconv.FormatUint(uint64(typed), 10)
	case uint8:
		return strconv.FormatUint(uint64(typed), 10)
	case uint16:
		return strconv.FormatUint(uint64(typed), 10)
	case uint32:
		return strconv.FormatUint(uint64(typed), 10)
	case uint64:
		return strconv.FormatUint(typed, 10)
	case float32:
		return strconv.FormatFloat(float64(typed), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return quoteSQLString(fmt.Sprint(value))
	}
}

// quoteSQLString 转义SQL字符串
func quoteSQLString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
