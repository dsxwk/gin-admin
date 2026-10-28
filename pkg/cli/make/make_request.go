package make

import (
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

type MakeRequest struct {
	base.Command
}

// NewMakeRequest 创建命令
func NewMakeRequest() *MakeRequest {
	return &MakeRequest{
		Command: base.NewCommand(base.Config{
			Name:        "make:request",
			Description: "验证请求创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: user",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "t", Long: "table"},
					Desc:     "表名, 如: roles",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "请求验证"},
					Desc:     "描述",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "c", Long: "camel", Default: "true"},
					Desc:     "字段是否使用驼峰",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "C", Long: "connection", Default: "mysql"},
					Desc:     "数据库连接",
					Required: false,
				},
			},
		}),
	}
}

// Handle 执行请求生成
func (m *MakeRequest) Handle(values map[string]string) error {

	file := m.GetMakeFile(values["file"], "request")
	templateName := "request"

	structName := lo.PascalCase(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))

	var fields []Field

	// table模式
	if values["table"] != "" {
		structName = lo.PascalCase(values["table"])
		camel := m.StringToBool(values["camel"])
		fields = m.loadTableFields(values["connection"], values["table"], camel)
	} else {
		// 默认ID模式
		fields = []Field{
			{
				Name:     "ID",
				JSON:     "id",
				Type:     "int64",
				Label:    "ID",
				Required: true,
				IsID:     true,
			},
		}
	}

	// 计算对齐宽度
	m.formatFields(fields)

	m.generateFile(templateName, file, structName, fields, values["desc"])
	return nil
}

type Field struct {
	FormattedField string // 格式化后的完整字段行
	Name           string
	JSON           string
	Type           string
	Label          string
	Required       bool
	Validation     string
	IsID           bool
}

// TableColumn 表字段结构
type TableColumn struct {
	Name     string // 字段名
	DataType string // 数据库类型
	Nullable bool   // 是否可为空
	Comment  string // 字段注释
}

// TranslateField 翻译字段
type TranslateField struct {
	Name  string
	Label string
}

// toJsonName 转换JSON字段名
func (m *MakeRequest) toJsonName(name string, camel bool) string {
	if !camel {
		return name
	}

	return lo.CamelCase(name)
}

// formatFields 预处理字段对齐
func (m *MakeRequest) formatFields(fields []Field) {
	// 计算最大宽度
	maxNameLen := 0
	maxTypeLen := 0
	for _, f := range fields {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
		if len(f.Type) > maxTypeLen {
			maxTypeLen = len(f.Type)
		}
	}

	// 生成格式化后的字段行
	for i := range fields {
		f := &fields[i]
		paddedName := fmt.Sprintf("%-*s", maxNameLen, f.Name)
		paddedType := fmt.Sprintf("%-*s", maxTypeLen, f.Type)
		f.FormattedField = fmt.Sprintf(
			"\t%s %s `json:\"%s\" form:\"%s\" label:\"%s\"`",
			paddedName,
			paddedType,
			f.JSON,
			f.JSON,
			f.Label,
		)
		f.Validation = m.buildFieldValidation(*f)
	}
}

// generateFile 生成请求文件
func (m *MakeRequest) generateFile(templateName, file, structName string, fields []Field, desc string) {
	templateFile := m.GetTemplate(templateName)

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("template parse error: %s", err.Error())
		os.Exit(1)
	}

	f := m.CheckDirAndFile(file)
	if f == nil {
		os.Exit(1)
	}

	packageName := filepath.Base(filepath.Dir(file))

	// 格式化输出翻译字段
	var translateFields []TranslateField
	for _, field := range fields {
		translateFields = append(translateFields, TranslateField{
			Name:  field.JSON,
			Label: field.Label,
		})
	}

	// 添加Page和PageSize到翻译字段列表计算最大宽度
	allTranslateFields := append([]TranslateField{}, translateFields...)
	allTranslateFields = append(allTranslateFields, TranslateField{Name: "page", Label: "页码"})
	allTranslateFields = append(allTranslateFields, TranslateField{Name: "pageSize", Label: "每页数量"})

	// 计算最大字段名长度用于对齐
	maxNameLen := 0
	for _, field := range allTranslateFields {
		if len(field.Name) > maxNameLen {
			maxNameLen = len(field.Name)
		}
	}

	// 生成格式化的翻译行
	var formattedTranslates []string
	for _, field := range translateFields {
		// 格式: "Name": "Label",
		line := fmt.Sprintf("\"%s\":%s%q,",
			field.Name,
			strings.Repeat(" ", maxNameLen-len(field.Name)+1),
			field.Label,
		)
		formattedTranslates = append(formattedTranslates, line)
	}

	// 添加Page
	formattedTranslates = append(formattedTranslates, fmt.Sprintf("\"page\":%s%q,",
		strings.Repeat(" ", maxNameLen-len("page")+1),
		"页码",
	))

	// 添加PageSize
	formattedTranslates = append(formattedTranslates, fmt.Sprintf("\"pageSize\":%s%q,",
		strings.Repeat(" ", maxNameLen-len("pageSize")+1),
		"每页数量",
	))

	idJSON := "id"
	for _, field := range fields {
		if field.IsID && field.JSON != "" {
			idJSON = field.JSON
			break
		}
	}

	data := struct {
		Package             string
		StructName          string
		Description         string
		Fields              []Field
		FormattedTranslates []string
		IDJSON              string
	}{
		Package:             packageName,
		StructName:          structName,
		Description:         desc,
		Fields:              fields,
		FormattedTranslates: formattedTranslates,
		IDJSON:              idJSON,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("template execute error: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("request generated: " + file)
}

// 加载表字段
func (m *MakeRequest) loadTableFields(conn, table string, camel bool) []Field {
	db := facade.DB(conn)

	// 获取表字段
	cols, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		terminal.Errorf("load table fields error: %s", err.Error())
		return []Field{}
	}

	var columns []TableColumn
	for _, col := range cols {
		name := col.Name()
		dt := strings.ToLower(col.DatabaseTypeName())
		nullable, _ := col.Nullable()
		comment, _ := col.Comment()

		columns = append(columns, TableColumn{
			Name:     name,
			DataType: dt,
			Nullable: nullable,
			Comment:  comment,
		})
	}

	if len(columns) == 0 {
		terminal.Errorf("No columns found for table: %s", table)
		return []Field{}
	}

	var fields []Field

	for _, col := range columns {
		// 过滤系统字段
		if col.Name == model.CreatedField ||
			col.Name == model.UpdatedField ||
			col.Name == model.DeletedField {
			continue
		}

		// 构建字段
		field := Field{
			Name:     m.toGoName(col.Name, camel),
			JSON:     m.toJsonName(col.Name, camel),
			Label:    m.parseComment(col.Comment, col.Name),
			Type:     m.getGoType(col.DataType),
			Required: !col.Nullable,
		}

		if col.Name == "id" {
			field.Required = true
			field.IsID = true
		}

		fields = append(fields, field)
	}

	return fields
}

// getGoType 将数据库类型转换为Go类型
func (m *MakeRequest) getGoType(dbType string) string {
	t := strings.ToLower(dbType)

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
	// 时间
	case strings.Contains(t, "time"),
		t == "date":
		return "string"
	// 二进制
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "[]byte"
	// json
	case t == "json", t == "jsonb":
		return "any"
	}

	return "string"
}

// buildFieldValidation 构建字段校验表达式
func (m *MakeRequest) buildFieldValidation(field Field) string {
	path := fmt.Sprintf("%q", field.JSON)
	value := "s." + field.Name

	switch field.Type {
	case "string":
		validation := fmt.Sprintf("fields.Str(%s, %s)", path, value)
		if field.Required {
			validation += ".Required()"
		}
		return validation + ".Errors(),"
	case "int64":
		validation := fmt.Sprintf("fields.Int(%s, %s)", path, value)
		if field.Required {
			validation += ".Required()"
		}
		return validation + ".Errors(),"
	case "float64":
		validation := fmt.Sprintf("fields.Float(%s, %s)", path, value)
		if field.Required {
			validation += ".Required()"
		}
		return validation + ".Errors(),"
	case "[]byte":
		validation := fmt.Sprintf("fields.Slice(%s, %s)", path, value)
		if field.Required {
			validation += ".Required()"
		}
		return validation + ".Errors(),"
	default:
		return ""
	}
}

// parseComment 解析字段注释
func (m *MakeRequest) parseComment(comment, fallback string) string {
	if comment == "" {
		return fallback
	}
	return comment
}

// toGoName 转换Go字段名
func (m *MakeRequest) toGoName(name string, camel bool) string {
	if !camel {
		return name
	}

	// 特殊处理id字段
	if name == "id" {
		return "ID"
	}

	// 转换下划线命名为大驼峰命名
	return lo.PascalCase(name)
}
