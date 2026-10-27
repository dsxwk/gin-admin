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
	"gorm.io/gorm"
)

// GrpcRequestField grpc请求字段
type GrpcRequestField struct {
	FormattedField       string
	FormattedUpdateField string
	FormattedTranslate   string
	Name                 string
	JSON                 string
	Type                 string
	Label                string
	Required             bool
	Validation           string
}

// MakeGrpcRequest grpc请求生成命令
type MakeGrpcRequest struct {
	base.Command
}

// NewMakeGrpcRequest 创建命令
func NewMakeGrpcRequest() *MakeGrpcRequest {
	return &MakeGrpcRequest{
		Command: base.NewCommand(base.Config{
			Name:        "grpc-make:request",
			Description: "grpc请求创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "t", Long: "table"},
					Desc:     "表名, 如: user 或 user,menu",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "p", Long: "path", Default: "grpc/request"},
					Desc:     "输出目录, 默认grpc/request",
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

// Handle 执行命令
func (m *MakeGrpcRequest) Handle(values map[string]string) error {
	outDir := strings.TrimPrefix(values["path"], "/")
	if outDir == "" {
		outDir = "grpc/request"
	}

	conn := values["connection"]
	db := facade.DB(conn)
	for _, table := range strings.Split(values["table"], ",") {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		terminal.Infof("开始生成grpc请求: %s", table)
		m.generateRequest(db, table, outDir)
	}
	return nil
}

// generateRequest 根据表结构生成grpc请求
func (m *MakeGrpcRequest) generateRequest(db *gorm.DB, table, outDir string) {
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		terminal.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	name := lo.PascalCase(table)
	var fields []GrpcRequestField
	for _, c := range columns {
		if c.Name == "id" ||
			c.Name == model.CreatedField ||
			c.Name == model.UpdatedField ||
			c.Name == model.DeletedField {
			continue
		}
		fieldName := grpcRequestFieldName(c.Name)
		jsonName := lo.CamelCase(fieldName)
		label := c.Comment
		if label == "" {
			label = c.Name
		}
		fields = append(fields, GrpcRequestField{
			Name:     fieldName,
			JSON:     jsonName,
			Type:     grpcRequestGoType(c.DataType),
			Label:    label,
			Required: !c.Nullable,
		})
	}

	m.formatGrpcRequestFields(fields)
	formattedId, formattedUpdateId := m.formatGrpcIdFields(fields)
	translates := m.formatGrpcRequestTranslates(fields)

	tpl, err := template.ParseFiles(m.GetTemplate("grpc_request"))
	if err != nil {
		terminal.Errorf("解析grpc请求模板失败: %s", err.Error())
		os.Exit(1)
	}

	data := struct {
		Package           string
		Name              string
		Var               string
		Description       string
		Fields            []GrpcRequestField
		FormattedId       string
		FormattedUpdateId string
		Translates        []string
	}{
		Package:           filepath.Base(outDir),
		Name:              name,
		Var:               lo.CamelCase(table),
		Description:       grpcRequestDescription(db, table),
		Fields:            fields,
		FormattedId:       formattedId,
		FormattedUpdateId: formattedUpdateId,
		Translates:        translates,
	}

	file := filepath.Join(outDir, table+".go")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	if err = tpl.Execute(f, data); err != nil {
		terminal.Errorf("生成grpc请求失败: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("grpc请求文件: " + file + " 生成成功!")
}

// formatGrpcIdFields 格式化grpc请求Id字段
func (m *MakeGrpcRequest) formatGrpcIdFields(fields []GrpcRequestField) (string, string) {
	maxNameLen := len("Id")
	maxTypeLen := len("int32")
	for _, f := range fields {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
		if len(f.Type) > maxTypeLen {
			maxTypeLen = len(f.Type)
		}
	}

	formattedId := "Id int32 `json:\"id\" label:\"ID\"`"
	formattedUpdateId := fmt.Sprintf(
		"%-*s %-*s `json:\"id\" label:\"ID\"`",
		maxNameLen,
		"Id",
		maxTypeLen,
		"int32",
	)
	return formattedId, formattedUpdateId
}

// formatGrpcRequestFields 格式化grpc请求字段
func (m *MakeGrpcRequest) formatGrpcRequestFields(fields []GrpcRequestField) {
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

	for i := range fields {
		f := &fields[i]
		paddedName := fmt.Sprintf("%-*s", maxNameLen, f.Name)
		paddedType := fmt.Sprintf("%-*s", maxTypeLen, f.Type)
		f.FormattedField = fmt.Sprintf(
			"%s %s `json:\"%s\" form:\"%s\" label:\"%s\"`",
			paddedName,
			paddedType,
			f.JSON,
			f.JSON,
			f.Label,
		)
		f.FormattedUpdateField = fmt.Sprintf(
			"%s %s `json:\"%s\" form:\"%s\" label:\"%s\"`",
			paddedName,
			paddedType,
			f.JSON,
			f.JSON,
			f.Label,
		)
		f.Validation = buildGrpcRequestValidation(*f)
	}
}

// formatGrpcRequestTranslates 格式化grpc请求翻译字段
func (m *MakeGrpcRequest) formatGrpcRequestTranslates(fields []GrpcRequestField) []string {
	type translateItem struct {
		name  string
		label string
	}
	items := []translateItem{
		{name: "Id", label: "ID"},
		{name: "Page", label: "页码"},
		{name: "PageSize", label: "每页数量"},
		{name: "Ids", label: "ID列表"},
	}
	for _, f := range fields {
		items = append(items, translateItem{name: f.Name, label: f.Label})
	}

	maxNameLen := 0
	for _, item := range items {
		if len(item.name) > maxNameLen {
			maxNameLen = len(item.name)
		}
	}

	var lines []string
	for _, item := range items {
		lines = append(lines, fmt.Sprintf(
			"\"%s\":%s%q,",
			item.name,
			strings.Repeat(" ", maxNameLen-len(item.name)+1),
			item.label,
		))
	}
	return lines
}

// grpcRequestDescription 获取grpc请求描述
func grpcRequestDescription(db *gorm.DB, table string) string {
	comment, err := getTableComment(db, table)
	if err != nil || comment == "" {
		return "请求"
	}
	return comment
}

// grpcRequestFieldName 数据库字段名转grpc请求Go字段名
func grpcRequestFieldName(name string) string {
	if name == "id" {
		return "Id"
	}
	return lo.PascalCase(name)
}

// grpcRequestGoType 数据库类型转grpc请求Go类型
func grpcRequestGoType(dbType string) string {
	t := strings.ToLower(dbType)

	switch {
	case strings.Contains(t, "int"):
		return "int32"
	case t == "bool" || t == "boolean":
		return "bool"
	case strings.Contains(t, "char"),
		strings.Contains(t, "text"),
		t == "uuid",
		t == "varchar":
		return "string"
	case strings.Contains(t, "float"),
		strings.Contains(t, "double"),
		strings.Contains(t, "decimal"),
		strings.Contains(t, "numeric"):
		return "float64"
	case t == "json", t == "jsonb":
		return "any"
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		return "string"
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "[]byte"
	}

	return "string"
}

// buildGrpcRequestValidation 构建grpc请求字段校验表达式
func buildGrpcRequestValidation(field GrpcRequestField) string {
	path := fmt.Sprintf("%q", field.Name)
	value := "s." + field.Name

	switch field.Type {
	case "string":
		validation := fmt.Sprintf("fields.Str(%s, %s)", path, value)
		if field.Required {
			validation += ".Required()"
		}
		return validation + ".Errors(),"
	case "int32":
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
