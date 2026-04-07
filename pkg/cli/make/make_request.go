package make

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeRequest struct {
	base.BaseCommand
}

func (m *MakeRequest) Name() string {
	return "make:request"
}

func (m *MakeRequest) Description() string {
	return "验证请求创建"
}

func (m *MakeRequest) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{Short: "f", Long: "file"},
			"文件路径, 如: user",
			true,
		},
		{
			base.Flag{Short: "t", Long: "table"},
			"表名, 如: roles",
			false,
		},
		{
			base.Flag{Short: "d", Long: "desc", Default: "请求验证"},
			"描述",
			false,
		},
		{
			base.Flag{Short: "c", Long: "camel", Default: "true"},
			"字段是否使用驼峰",
			false,
		},
		{
			base.Flag{Short: "C", Long: "connection", Default: "mysql"},
			"数据库连接",
			false,
		},
	}
}

func init() {
	cli.Register(&MakeRequest{})
}

type Field struct {
	Name     string
	JSON     string
	Type     string
	Label    string
	Validate string
	IsID     bool
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

func (m *MakeRequest) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())

	file := m.GetMakeFile(values["file"], "request")
	templateName := "request"

	structName := pkg.SnakeToCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))

	var fields []Field

	// table模式
	if values["table"] != "" {
		structName = pkg.SnakeToCamel(values["table"])
		fields = m.loadTableFields(values["connection"], values["table"], values["camel"] == "true")
	} else {
		// 默认ID模式
		fields = []Field{
			{
				Name:     "ID",
				JSON:     "id",
				Type:     "int64",
				Label:    "ID",
				Validate: "required|int|gt:0",
				IsID:     true,
			},
		}
	}

	m.generateFile(templateName, file, structName, fields, values["desc"])
}

func (m *MakeRequest) generateFile(templateName, file, structName string, fields []Field, desc string) {
	templateFile := m.GetTemplate(templateName)

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("template parse error: %s", err.Error())
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
			Name:  field.Name,
			Label: field.Label,
		})
	}

	// 添加Page和PageSize到翻译字段列表计算最大宽度
	allTranslateFields := append([]TranslateField{}, translateFields...)
	allTranslateFields = append(allTranslateFields, TranslateField{Name: "Page", Label: "页码"})
	allTranslateFields = append(allTranslateFields, TranslateField{Name: "PageSize", Label: "每页数量"})

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
			strings.Repeat(" ", maxNameLen-len(field.Name)+2),
			field.Label,
		)
		formattedTranslates = append(formattedTranslates, line)
	}

	// 添加Page
	formattedTranslates = append(formattedTranslates, fmt.Sprintf("\"Page\":%s%q,",
		strings.Repeat(" ", maxNameLen-len("Page")+2),
		"页码",
	))

	// 添加PageSize
	formattedTranslates = append(formattedTranslates, fmt.Sprintf("\"PageSize\":%s%q,",
		strings.Repeat(" ", maxNameLen-len("PageSize")+2),
		"每页数量",
	))

	data := struct {
		Package             string
		StructName          string
		Description         string
		Fields              []Field
		CreateScene         string
		UpdateScene         string
		FormattedTranslates []string
	}{
		Package:             packageName,
		StructName:          structName,
		Description:         desc,
		Fields:              fields,
		CreateScene:         m.buildScene(fields, false),
		UpdateScene:         m.buildScene(fields, true),
		FormattedTranslates: formattedTranslates,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("template execute error: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("request generated: " + file)
}

// 加载表字段
func (m *MakeRequest) loadTableFields(conn, table string, camel bool) []Field {
	db := facade.DB.Connection(conn)

	// 获取表字段
	cols, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		flag.Errorf("load table fields error: %s", err.Error())
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
		flag.Errorf("No columns found for table: %s", table)
		return []Field{}
	}

	var fields []Field

	for _, col := range columns {
		// 过滤系统字段
		if col.Name == "created_at" ||
			col.Name == "updated_at" ||
			col.Name == "deleted_at" {
			continue
		}

		// 构建字段
		field := Field{
			Name:  m.toGoName(col.Name, camel),
			JSON:  col.Name,
			Label: m.parseComment(col.Comment, col.Name),
			Type:  m.getGoType(col.DataType),
		}

		// 构建验证规则
		if col.Name == "id" {
			field.Validate = "required|int|gt:0"
			field.IsID = true
		} else {
			rules := []string{"required"}
			rules = append(rules, m.getValidateRules(col.DataType)...)
			field.Validate = strings.Join(rules, "|")
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
		return "string"
	}

	return "string"
}

// getValidateRules 根据数据库类型获取验证规则
func (m *MakeRequest) getValidateRules(dbType string) []string {
	t := strings.ToLower(dbType)
	var rules []string

	// 整数类型添加int验证
	if strings.Contains(t, "int") {
		rules = append(rules, "int")
	}

	// 字符串类型添加max验证
	if strings.Contains(t, "char") || strings.Contains(t, "text") {
		rules = append(rules, "max:255")
	}

	return rules
}

// buildScene 生成字符串数组
func (m *MakeRequest) buildScene(fields []Field, update bool) string {
	var arr []string

	for _, f := range fields {
		if update && f.IsID {
			arr = append(arr, `"ID"`)
			continue
		}
		if !update && f.IsID {
			continue
		}
		arr = append(arr, `"`+f.Name+`"`)
	}

	return strings.Join(arr, ", ")
}

func (m *MakeRequest) parseComment(comment, fallback string) string {
	if comment == "" {
		return fallback
	}
	return comment
}

func (m *MakeRequest) toGoName(name string, camel bool) string {
	if !camel {
		return name
	}

	// 特殊处理id字段
	if name == "id" {
		return "ID"
	}

	// 转换下划线命名为大驼峰命名
	return pkg.ToUpperCamel(name)
}
