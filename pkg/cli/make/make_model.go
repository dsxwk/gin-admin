package make

import (
	"bytes"
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/fatih/color"
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
	Name     string // 字段名
	DataType string // 数据库类型
	Nullable bool   // 是否可为空
	Comment  string // 字段注释
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
			base.Flag{
				Short: "t",
				Long:  "table",
			},
			"表名, 如: user 或 user,menu",
			true,
		},
		{
			base.Flag{
				Short: "p",
				Long:  "path",
			},
			"输出目录, 如: api/user",
			false,
		},
		{
			base.Flag{
				Short:   "c",
				Long:    "camel",
				Default: "true",
			},
			"json字段是否使用驼峰",
			false,
		},
		{
			base.Flag{
				Short:   "C",
				Long:    "connection",
				Default: "mysql",
			},
			"数据库连接",
			false,
		},
	}
}

// Execute cli命令执行入口
func (m *MakeModel) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
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
	db := facade.DB.Connection(conn)
	for _, table := range tables {
		color.Cyan("开始生成模型: %s", table)

		err := m.generateModel(_make, db, table, path, camel)
		if err != nil {
			flag.Errorf("生成失败: %s", err.Error())
			continue
		}

		flag.Successf("模型生成成功: " + filepath.Join(path, table+".go"))
	}
}

// init 注册cli命令
func init() {
	cli.Register(&MakeModel{})
}

// generateModel 根据表结构生成 Model 文件
func (m *MakeModel) generateModel(_make string, db *gorm.DB, table string, outDir string, camel bool) error {
	// 获取表字段
	cols, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		return err
	}

	im := NewImport()

	// 当前包名
	pkgName := filepath.Base(outDir)

	var columns []Column
	for _, col := range cols {
		name := col.Name()
		dt := strings.ToLower(col.DatabaseTypeName())
		nullable, _ := col.Nullable()
		comment, _ := col.Comment()

		columns = append(columns, Column{
			Name:     name,
			DataType: dt,
			Nullable: nullable,
			Comment:  comment,
		})
	}

	structName := snakeToCamel(table)
	tableConst := "TableName" + structName

	// 计算字段对齐长度
	maxNameLen := 0
	maxTypeLen := 0

	for _, c := range columns {
		name := snakeToCamel(c.Name)
		typ := goType(c, im, pkgName)

		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}

		if len(typ) > maxTypeLen {
			maxTypeLen = len(typ)
		}
	}

	var fieldLines []string

	for _, c := range columns {
		fieldName := snakeToCamel(c.Name)
		fieldType := goType(c, im, pkgName)

		var jsonName string
		if camel {
			jsonName = snakeToLowerCamel(c.Name)
		} else {
			jsonName = c.Name
		}

		tag := fmt.Sprintf(
			"`%s json:\"%s\" form:\"%s\"`",
			buildGormTag(c),
			jsonName,
			jsonName,
		)

		// deleted_at自动忽略swagger
		if jsonName == "deletedAt" || c.Name == "deleted_at" {
			tag = strings.TrimSuffix(tag, "`") + " swaggerignore:\"true\"`"
		}

		line := fmt.Sprintf("%-*s %-*s %s", maxNameLen, fieldName, maxTypeLen, fieldType, tag)

		fieldLines = append(fieldLines, line)
	}

	templateFile := m.GetTemplate(_make)
	tpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}
	data := struct {
		Imports    string
		Struct     string
		Table      string
		TableConst string
		Fields     []string
	}{
		Imports:    im.Render(),
		Struct:     structName,
		Table:      table,
		TableConst: tableConst,
		Fields:     fieldLines,
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		return err
	}

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}

	file := filepath.Join(outDir, table+".go")

	return os.WriteFile(file, buf.Bytes(), 0644)
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
		t == "uuid":
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
	case strings.Contains(t, "time"),
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

// buildGormTag 生成gorm tag
func buildGormTag(c Column) string {
	var tags []string

	tags = append(tags, "column:"+c.Name)

	if c.Comment != "" {
		tags = append(tags, "comment:"+c.Comment)
	}

	return "gorm:\"" + strings.Join(tags, ";") + "\""
}

// snakeToCamel 下划线转大驼峰
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")

	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}

	return strings.Join(parts, "")
}

// snakeToLowerCamel 下划线转小驼峰
func snakeToLowerCamel(s string) string {
	camel := snakeToCamel(s)

	return strings.ToLower(camel[:1]) + camel[1:]
}
