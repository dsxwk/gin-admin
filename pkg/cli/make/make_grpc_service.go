package make

import (
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// GrpcServiceField grpc服务字段
type GrpcServiceField struct {
	Name   string
	JSON   string
	Scalar bool
	IsJSON bool
}

// MakeGrpcService grpc服务生成命令
type MakeGrpcService struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeGrpcService) Name() string {
	return "grpc-make:service"
}

// Description 命令描述
func (m *MakeGrpcService) Description() string {
	return "grpc服务创建"
}

// Help 命令选项
func (m *MakeGrpcService) Help() []base.CommandOption {
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
				Default: "grpc/service",
			},
			"输出目录, 默认grpc/service",
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
				Short:   "a",
				Long:    "auth",
				Default: "true",
			},
			"是否需要鉴权, 默认true",
			false,
		},
	}
}

// Execute 执行命令
func (m *MakeGrpcService) Execute(values map[string]string) {
	outDir := strings.TrimPrefix(values["path"], "/")
	if outDir == "" {
		outDir = "grpc/service"
	}

	conn := values["connection"]
	auth := m.StringToBool(values["auth"])
	db := facade.DB(conn)
	for _, table := range strings.Split(values["table"], ",") {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		flag.Infof("开始生成grpc服务: %s", table)
		m.generateService(db, table, outDir, auth)
	}
}

// generateService 根据表结构生成grpc服务
func (m *MakeGrpcService) generateService(db *gorm.DB, table, outDir string, auth bool) {
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		flag.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	name := lo.PascalCase(table)
	var hasCreated, hasUpdated bool
	var fields []GrpcServiceField
	for _, c := range columns {
		if c.Name == model.CreatedField {
			hasCreated = true
		}
		if c.Name == model.UpdatedField {
			hasUpdated = true
		}
		if c.Name == "id" ||
			c.Name == model.CreatedField ||
			c.Name == model.UpdatedField ||
			c.Name == model.DeletedField {
			continue
		}
		fieldName := lo.PascalCase(c.Name)
		fields = append(fields, GrpcServiceField{
			Name:   fieldName,
			JSON:   lo.CamelCase(fieldName),
			Scalar: grpcServiceScalar(c.DataType),
			IsJSON: grpcServiceIsJSON(c.DataType),
		})
	}

	tpl, err := template.ParseFiles(m.GetTemplate("grpc_service"))
	if err != nil {
		flag.Errorf("解析grpc服务模板失败: %s", err.Error())
		os.Exit(1)
	}

	data := struct {
		Package             string
		Name                string
		Var                 string
		Description         string
		CreateFields        string
		ProtoFields         string
		RequestCreateFields string
		UpdateJsonFields    string
		Auth                bool
	}{
		Package:             filepath.Base(outDir),
		Name:                name,
		Var:                 lo.CamelCase(table),
		Description:         grpcServiceDescription(db, table),
		CreateFields:        buildGrpcCreateFields(fields),
		ProtoFields:         buildGrpcModelProtoFields(fields, hasCreated, hasUpdated, lo.CamelCase(table)),
		RequestCreateFields: buildGrpcRequestFields(fields, lo.CamelCase(table)),
		UpdateJsonFields:    buildGrpcUpdateJsonFields(fields),
		Auth:                auth,
	}

	file := filepath.Join(outDir, table+".go")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	if err = tpl.Execute(f, data); err != nil {
		flag.Errorf("生成grpc服务失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("grpc服务文件: " + file + " 生成成功!")
}

// grpcServiceDescription 获取grpc服务描述
func grpcServiceDescription(db *gorm.DB, table string) string {
	comment, err := getTableComment(db, table)
	if err != nil || comment == "" {
		return "grpc服务"
	}
	return comment
}

// grpcServiceScalar 判断是否为可自动转换的标量字段
func grpcServiceScalar(dbType string) bool {
	t := strings.ToLower(dbType)
	switch {
	case strings.Contains(t, "int"),
		t == "bool",
		t == "boolean",
		strings.Contains(t, "char"),
		strings.Contains(t, "text"),
		t == "uuid",
		t == "varchar",
		strings.Contains(t, "float"),
		strings.Contains(t, "double"),
		strings.Contains(t, "decimal"),
		strings.Contains(t, "numeric"):
		return true
	}
	return false
}

// grpcServiceIsJSON 判断是否为JSON字段
func grpcServiceIsJSON(dbType string) bool {
	t := strings.ToLower(dbType)
	return t == "json" || t == "jsonb"
}

// buildGrpcCreateFields 构建创建模型字段
func buildGrpcCreateFields(fields []GrpcServiceField) string {
	if len(fields) == 0 {
		return ""
	}

	maxLen := 0
	for _, f := range fields {
		if (f.Scalar || f.IsJSON) && len(f.Name) > maxLen {
			maxLen = len(f.Name)
		}
	}

	var lines []string
	for _, f := range fields {
		switch {
		case !f.Scalar && !f.IsJSON:
			continue
		case f.IsJSON:
			lines = append(lines, fmt.Sprintf(
				"        %s:%s &model.JsonValue{Data: dto.%s},",
				f.Name,
				strings.Repeat(" ", maxLen-len(f.Name)+1),
				f.Name,
			))
		default:
			lines = append(lines, fmt.Sprintf(
				"        %s:%s dto.%s,",
				f.Name,
				strings.Repeat(" ", maxLen-len(f.Name)+1),
				f.Name,
			))
		}
	}
	return strings.Join(lines, "\n")
}

// buildGrpcModelProtoFields 构建模型转proto字段
func buildGrpcModelProtoFields(fields []GrpcServiceField, hasCreated, hasUpdated bool, varName string) string {
	type protoItem struct {
		name  string
		value string
	}
	items := []protoItem{{name: "Id", value: "m.ID"}}
	for _, f := range fields {
		switch {
		case !f.Scalar && !f.IsJSON:
			continue
		case f.IsJSON:
			items = append(items, protoItem{name: f.Name, value: varName + "ToStruct(m." + f.Name + ")"})
		default:
			items = append(items, protoItem{name: f.Name, value: "m." + f.Name})
		}
	}
	if hasCreated {
		items = append(items, protoItem{name: "CreatedAt", value: "m.CreatedAt.String()"})
	}
	if hasUpdated {
		items = append(items, protoItem{name: "UpdatedAt", value: "m.UpdatedAt.String()"})
	}

	maxLen := 0
	for _, item := range items {
		if len(item.name) > maxLen {
			maxLen = len(item.name)
		}
	}

	var lines []string
	for _, item := range items {
		lines = append(lines, fmt.Sprintf(
			"        %s:%s %s,",
			item.name,
			strings.Repeat(" ", maxLen-len(item.name)+1),
			item.value,
		))
	}
	return strings.Join(lines, "\n")
}

// buildGrpcRequestFields 构建proto请求转grpc请求字段
func buildGrpcRequestFields(fields []GrpcServiceField, varName string) string {
	if len(fields) == 0 {
		return ""
	}

	maxLen := 0
	for _, f := range fields {
		if (f.Scalar || f.IsJSON) && len(f.Name) > maxLen {
			maxLen = len(f.Name)
		}
	}

	var lines []string
	for _, f := range fields {
		switch {
		case !f.Scalar && !f.IsJSON:
			continue
		case f.IsJSON:
			lines = append(lines, fmt.Sprintf(
				"        %s:%s %sStructToMap(req.%s),",
				f.Name,
				strings.Repeat(" ", maxLen-len(f.Name)+1),
				varName,
				f.Name,
			))
		default:
			lines = append(lines, fmt.Sprintf(
				"        %s:%s req.%s,",
				f.Name,
				strings.Repeat(" ", maxLen-len(f.Name)+1),
				f.Name,
			))
		}
	}
	return strings.Join(lines, "\n")
}

// buildGrpcUpdateJsonFields 构建更新map的JSON字段转换
func buildGrpcUpdateJsonFields(fields []GrpcServiceField) string {
	var lines []string
	for _, f := range fields {
		if !f.IsJSON {
			continue
		}
		jsonName := lo.SnakeCase(f.Name)
		lines = append(lines, fmt.Sprintf("    if v, ok := rows[%q]; ok {", jsonName))
		lines = append(lines, fmt.Sprintf("        rows[%q] = &model.JsonValue{Data: v}", jsonName))
		lines = append(lines, "    }")
	}
	return strings.Join(lines, "\n")
}

func init() {
	cli.Register(&MakeGrpcService{})
}
