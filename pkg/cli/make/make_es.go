package make

import (
	"bytes"
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// MakeEs ES搜索生成命令
type MakeEs struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeEs) Name() string {
	return "make:es"
}

// Description 命令描述
func (m *MakeEs) Description() string {
	return "ES搜索创建"
}

// Help 命令选项
func (m *MakeEs) Help() []base.CommandOption {
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
				Short:   "p",
				Long:    "path",
				Default: "app/es",
			},
			"输出目录, 默认app/es",
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
		{
			base.Flag{
				Short:   "e",
				Long:    "exclude",
				Default: "password",
			},
			"排除字段,多个用逗号分隔,如: password,token",
			false,
		},
	}
}

// Execute 执行命令
func (m *MakeEs) Execute(values map[string]string) {
	outDir := strings.TrimPrefix(strings.TrimSpace(values["path"]), "/")
	if outDir == "" {
		outDir = "app/es"
	}

	conn := values["connection"]
	exclude := make(map[string]bool)
	for name := range strings.SplitSeq(values["exclude"], ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			exclude[name] = true
		}
	}
	exclude[model.DeletedField] = true

	db := facade.DB(conn)
	for table := range strings.SplitSeq(values["table"], ",") {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		flag.Infof("开始生成ES搜索: %s", table)
		m.generateEs(db, table, outDir, exclude)
	}
}

// generateEs 根据表结构生成ES搜索
func (m *MakeEs) generateEs(db *gorm.DB, table, outDir string, exclude map[string]bool) {
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		flag.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	description, _ := getTableComment(db, table)
	description = cleanDescription(description, "ES搜索")
	name := lo.PascalCase(table)
	varName := lo.CamelCase(table)

	var fields []esField
	for _, column := range columns {
		if exclude[column.Name] {
			continue
		}
		field := esField{
			JSON:       lo.CamelCase(esFieldName(column.Name)),
			Type:       esMappingType(column),
			DocValue:   esDocValue(esFieldName(column.Name), column),
			TextSearch: esTextField(column),
		}
		fields = append(fields, field)
	}

	if len(fields) == 0 {
		flag.Errorf("表 %s 没有可生成的ES字段", table)
		os.Exit(1)
	}

	tpl, err := template.ParseFiles(m.GetTemplate("es"))
	if err != nil {
		flag.Errorf("解析ES搜索模板失败: %s", err.Error())
		os.Exit(1)
	}

	data := struct {
		Package       string
		Table         string
		Name          string
		Var           string
		Description   string
		MappingFields string
		UpdateFields  string
		TextFields    string
		DocFields     string
	}{
		Package:       filepath.Base(outDir),
		Table:         table,
		Name:          name,
		Var:           varName,
		Description:   description,
		MappingFields: m.buildMappingFields(fields, "\t\t"),
		UpdateFields:  m.buildUpdateFields(fields, "\t"),
		TextFields:    m.buildTextFields(fields, "\t"),
		DocFields:     m.buildDocFields(fields, "\t\t"),
	}

	file := filepath.Join(outDir, table+".go")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		flag.Errorf("生成ES搜索文件失败: %s", err.Error())
		os.Exit(1)
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		flag.Errorf("格式化ES搜索文件失败: %s", err.Error())
		os.Exit(1)
	}
	if _, err = f.Write(source); err != nil {
		flag.Errorf("写入ES搜索文件失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("ES搜索文件: " + file + " 生成成功!")
}

// esField ES字段信息
type esField struct {
	JSON       string // JSON字段名
	Type       string // ES映射类型
	DocValue   string // 文档值表达式
	TextSearch bool   // 是否参与关键字搜索
}

// esFieldName 数据库字段名转Go字段名
func esFieldName(name string) string {
	if name == "id" {
		return "ID"
	}
	return lo.PascalCase(name)
}

// esMappingType 数据库类型转ES映射类型
func esMappingType(c Column) string {
	t := strings.ToLower(c.DataType)
	switch {
	case strings.Contains(t, "bigint"):
		return "long"
	case strings.Contains(t, "tinyint"):
		return "byte"
	case strings.Contains(t, "smallint"):
		return "short"
	case strings.Contains(t, "int"):
		return "integer"
	case t == "bool" || t == "boolean":
		return "boolean"
	case strings.Contains(t, "float"):
		return "float"
	case strings.Contains(t, "double"),
		strings.Contains(t, "decimal"),
		strings.Contains(t, "numeric"):
		return "double"
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		return "date"
	case t == "json", t == "jsonb":
		return "object"
	case strings.Contains(t, "text"):
		return "text"
	case strings.Contains(t, "char"),
		strings.Contains(t, "varchar"),
		t == "uuid",
		t == "enum",
		t == "set":
		return "keyword"
	default:
		return "keyword"
	}
}

// esTextField 是否为文本搜索字段
func esTextField(c Column) bool {
	t := strings.ToLower(c.DataType)
	return strings.Contains(t, "char") ||
		strings.Contains(t, "text") ||
		t == "uuid" ||
		t == "enum" ||
		t == "set"
}

// esDocValue 获取文档值表达式
func esDocValue(name string, c Column) string {
	t := strings.ToLower(c.DataType)
	switch {
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		return "s.Client().FormatDateTime(m." + name + ")"
	case t == "json", t == "jsonb":
		return "s.Client().JsonValue(m." + name + ")"
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "string(m." + name + ")"
	default:
		return "m." + name
	}
}

// cleanDescription 清理表注释
func cleanDescription(comment, def string) string {
	comment = strings.TrimSpace(comment)
	comment = strings.NewReplacer("\n", " ", "\r", "", "\"", "").Replace(comment)
	if comment == "" {
		return def
	}
	return comment
}

// buildMappingFields 构建映射字段
func (m *MakeEs) buildMappingFields(fields []esField, indent string) string {
	maxLen := 0
	for _, field := range fields {
		if len(field.JSON) > maxLen {
			maxLen = len(field.JSON)
		}
	}

	lines := make([]string, 0, len(fields))
	for _, field := range fields {
		lines = append(lines, fmt.Sprintf(
			"%s%q:%s%q,",
			indent,
			field.JSON,
			strings.Repeat(" ", maxLen-len(field.JSON)+1),
			field.Type,
		))
	}
	return strings.Join(lines, "\n")
}

// buildUpdateFields 构建更新字段
func (m *MakeEs) buildUpdateFields(fields []esField, indent string) string {
	lines := make([]string, 0, len(fields))
	maxLen := 0
	for _, field := range fields {
		if field.JSON == "id" {
			continue
		}
		if len(field.JSON) > maxLen {
			maxLen = len(field.JSON)
		}
	}

	for _, field := range fields {
		if field.JSON == "id" {
			continue
		}
		lines = append(lines, fmt.Sprintf(
			"%s%q:%s true,",
			indent,
			field.JSON,
			strings.Repeat(" ", maxLen-len(field.JSON)+1),
		))
	}
	return strings.Join(lines, "\n")
}

// buildTextFields 构建文本搜索字段
func (m *MakeEs) buildTextFields(fields []esField, indent string) string {
	var lines []string
	for _, field := range fields {
		if !field.TextSearch {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s%q,", indent, field.JSON))
	}
	return strings.Join(lines, "\n")
}

// buildDocFields 构建文档字段
func (m *MakeEs) buildDocFields(fields []esField, indent string) string {
	maxLen := 0
	for _, field := range fields {
		if len(field.JSON) > maxLen {
			maxLen = len(field.JSON)
		}
	}

	lines := make([]string, 0, len(fields))
	for _, field := range fields {
		lines = append(lines, fmt.Sprintf(
			"%s%q:%s %s,",
			indent,
			field.JSON,
			strings.Repeat(" ", maxLen-len(field.JSON)+1),
			field.DocValue,
		))
	}
	return strings.Join(lines, "\n")
}

func init() {
	cli.Register(&MakeEs{})
}
