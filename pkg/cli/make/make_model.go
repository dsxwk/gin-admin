package make

import (
	"database/sql"
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// MakeModel 模型生成命令
type MakeModel struct {
	base.BaseCommand
}

// Column 表字段结构
type Column struct {
	Name         string         // 字段名
	DataType     string         // 数据库类型
	Nullable     bool           // 是否可为空
	Comment      string         // 字段注释
	ColumnType   string         // 完整列类型(如 tinyint(3) unsigned)
	DefaultValue sql.NullString // 默认值
	IsUnsigned   bool           // 是否无符号
	IsAutoIncr   bool           // 是否自增
	Length       int            // 字段长度
}

// Import 用于管理自动生成的import包
type Import struct {
	pkgs map[string]struct{}
}

// NewImport 创建import管理器
func NewImport() *Import {
	return &Import{
		pkgs: make(map[string]struct{}),
	}
}

// Add 添加import
func (m *Import) Add(pkg string) {
	if pkg != "" {
		m.pkgs[pkg] = struct{}{}
	}
}

// Render 渲染import代码
func (m *Import) Render() string {
	if len(m.pkgs) == 0 {
		return ""
	}

	var list []string
	for p := range m.pkgs {
		list = append(list, fmt.Sprintf("\t%q", p))
	}

	sort.Strings(list)

	return "import (\n" + strings.Join(list, "\n") + "\n)\n"
}

// Name 返回cli命令名称
func (m *MakeModel) Name() string {
	return "make:model"
}

// Description 命令描述
func (m *MakeModel) Description() string {
	return "模型创建"
}

// Help 返回命令参数说明
func (m *MakeModel) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			Flag: base.Flag{
				Short: "t",
				Long:  "table",
			},
			Desc:     "表名, 如: user 或 user,menu",
			Required: true,
		},
		{
			Flag: base.Flag{
				Short: "p",
				Long:  "path",
			},
			Desc:     "输出目录, 如: api/user",
			Required: false,
		},
		{
			Flag: base.Flag{
				Short:   "c",
				Long:    "camel",
				Default: "true",
			},
			Desc:     "json字段是否使用驼峰",
			Required: false,
		},
		{
			Flag: base.Flag{
				Short:   "C",
				Long:    "connection",
				Default: "mysql",
			},
			Desc:     "数据库连接",
			Required: false,
		},
	}
}

// Execute cli命令执行入口
func (m *MakeModel) Execute(values map[string]string) {
	// 输出目录
	path := filepath.Join("app/model/", strings.TrimPrefix(values["path"], "/"))
	// 表名支持多个
	tables := strings.Split(values["table"], ",")
	for i := range tables {
		tables[i] = strings.TrimSpace(tables[i])
	}

	camel := m.StringToBool(values["camel"])
	conn := values["connection"]
	_make := strings.TrimPrefix(m.Name(), "make:")
	db := facade.DB(conn)
	for _, table := range tables {
		flag.Infof("开始生成模型: %s", table)
		m.generateModel(_make, db, table, path, conn, camel)
	}
}

func init() {
	cli.Register(&MakeModel{})
}

// GetColumnInfo 获取完整字段信息
func GetColumnInfo(db *gorm.DB, tableName string) ([]Column, error) {
	var columns []Column

	// 使用raw SQL获取详细信息
	rows, err := db.Raw(`
        SELECT 
            COLUMN_NAME,
            COLUMN_TYPE,
            IS_NULLABLE,
            COLUMN_COMMENT,
            COLUMN_DEFAULT,
            EXTRA,
            DATA_TYPE,
            NUMERIC_PRECISION,
            CHARACTER_MAXIMUM_LENGTH
        FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = ?
        ORDER BY ORDINAL_POSITION
    `, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var col Column
		var isNullable string
		var extra string
		var numericPrecision, charMaxLength sql.NullInt64

		err = rows.Scan(
			&col.Name,
			&col.ColumnType,
			&isNullable,
			&col.Comment,
			&col.DefaultValue,
			&extra,
			&col.DataType,
			&numericPrecision,
			&charMaxLength,
		)
		if err != nil {
			return nil, err
		}

		col.Nullable = isNullable == "YES"
		col.IsAutoIncr = strings.Contains(extra, "auto_increment")
		col.IsUnsigned = strings.Contains(col.ColumnType, "unsigned")

		// 解析长度
		if numericPrecision.Valid {
			col.Length = int(numericPrecision.Int64)
		} else if charMaxLength.Valid {
			col.Length = int(charMaxLength.Int64)
		}

		columns = append(columns, col)
	}

	return columns, nil
}

// generateModel 根据表结构生成Model文件
func (m *MakeModel) generateModel(_make string, db *gorm.DB, table string, outDir, conn string, camel bool) {
	// 获取完整字段信息
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		flag.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	// 获取表注释
	tableComment, _ := getTableComment(db, table)

	im := NewImport()
	pkgName := filepath.Base(outDir)
	structName := lo.PascalCase(table)
	tableConst := "TableName" + structName

	// 计算字段对齐长度
	maxNameLen := 0
	maxTypeLen := 0
	maxTagLen := 0

	for _, c := range columns {
		name := lo.PascalCase(c.Name)
		typ := goType(c, im, pkgName)
		tag := buildGormTag(c)

		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
		if len(typ) > maxTypeLen {
			maxTypeLen = len(typ)
		}
		if len(tag) > maxTagLen {
			maxTagLen = len(tag)
		}
	}

	var fieldLines []string

	for _, c := range columns {
		var fieldName string
		if c.Name == "id" {
			fieldName = "ID"
		} else {
			fieldName = lo.PascalCase(c.Name)
		}

		fieldType := goType(c, im, pkgName)

		var jsonName string
		if camel {
			jsonName = lo.CamelCase(fieldName)
		} else {
			jsonName = c.Name
		}

		gormTag := buildGormTag(c)
		tag := fmt.Sprintf("`%s json:\"%s\" form:\"%s\"`", gormTag, jsonName, jsonName)

		// deleted_at 自动忽略 swagger
		if jsonName == "deletedAt" || c.Name == "deleted_at" {
			tag = strings.TrimSuffix(tag, "`") + " swaggerignore:\"true\"`"
		}

		// 添加字段注释
		//comment := ""
		//if c.Comment != "" {
		//	comment = fmt.Sprintf("// %s", c.Comment)
		//}
		//
		//line := fmt.Sprintf("%-*s %-*s %-*s %s", maxNameLen, fieldName, maxTypeLen, fieldType, maxTagLen, tag, comment)
		//fieldLines = append(fieldLines, strings.TrimSpace(line))

		// 不生成行尾注释
		line := fmt.Sprintf("%-*s %-*s %s", maxNameLen, fieldName, maxTypeLen, fieldType, tag)
		fieldLines = append(fieldLines, strings.TrimSpace(line))
	}

	templateFile := m.GetTemplate(_make)
	tpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}

	data := struct {
		Imports       string
		Struct        string
		StructComment string
		Table         string
		TableConst    string
		Connection    string
		Fields        []string
	}{
		Imports:       im.Render(),
		Struct:        structName,
		StructComment: tableComment,
		Table:         table,
		TableConst:    tableConst,
		Connection:    conn,
		Fields:        fieldLines,
	}

	file := filepath.Join(outDir, table+".go")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	err = tpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("模型文件: " + file + " 生成成功!")
}

// getTableComment 获取表注释
func getTableComment(db *gorm.DB, tableName string) (string, error) {
	var comment string
	err := db.Raw(`
        SELECT TABLE_COMMENT 
        FROM INFORMATION_SCHEMA.TABLES 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = ?
    `, tableName).Scan(&comment).Error
	if err != nil {
		return "", err
	}
	return comment, nil
}

// goType 将数据库类型转换为go类型
func goType(c Column, im *Import, pkgName string) string {
	t := strings.ToLower(c.DataType)

	switch {
	// 整数
	case strings.Contains(t, "int"):
		return "int64"

	// 布尔
	case t == "bool" || t == "boolean":
		return "bool"

	// 字符串
	case strings.Contains(t, "char"),
		strings.Contains(t, "text"),
		t == "uuid",
		t == "varchar":
		return "string"

	// 浮点
	case strings.Contains(t, "float"),
		strings.Contains(t, "double"),
		strings.Contains(t, "decimal"),
		strings.Contains(t, "numeric"):
		return "float64"

	// json
	case t == "json", t == "jsonb":
		if pkgName != "model" {
			im.Add("gin/app/model")
			return "*model.JsonValue"
		}
		return "*JsonValue"

	// 时间
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		if c.Name == "deleted_at" {
			if pkgName != "model" {
				im.Add("gin/app/model")
				return "*model.DeletedAt"
			}
			return "*DeletedAt"
		}
		if pkgName != "model" {
			im.Add("gin/app/model")
			return "*model.DateTime"
		}
		return "*DateTime"

	// 二进制
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "[]byte"
	}

	return "string"
}

// buildGormTag 生成完整的gorm tag
func buildGormTag(c Column) string {
	var tags []string

	tags = append(tags, "column:"+c.Name)

	// 主键自增
	if c.IsAutoIncr {
		tags = append(tags, "primaryKey;autoIncrement")
	}

	// 非空
	if !c.Nullable {
		tags = append(tags, "not null")
	}

	// 默认值
	if c.DefaultValue.Valid && c.DefaultValue.String != "" && c.DefaultValue.String != "NULL" {
		tags = append(tags, "default:"+c.DefaultValue.String)
	}

	// 类型
	if c.ColumnType != "" {
		tags = append(tags, "type:"+c.ColumnType)
	}

	// 注释
	if c.Comment != "" {
		tags = append(tags, "comment:"+c.Comment)
	}

	return "gorm:\"" + strings.Join(tags, ";") + "\""
}
