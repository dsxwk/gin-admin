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

// MakeGrpcProto grpc proto生成命令
type MakeGrpcProto struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeGrpcProto) Name() string {
	return "grpc-make:proto"
}

// Description 命令描述
func (m *MakeGrpcProto) Description() string {
	return "grpc proto创建"
}

// Help 命令选项
func (m *MakeGrpcProto) Help() []base.CommandOption {
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
				Default: "grpc/proto",
			},
			"输出目录, 默认grpc/proto",
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

// Execute 执行命令
func (m *MakeGrpcProto) Execute(values map[string]string) {
	outDir := strings.TrimPrefix(values["path"], "/")
	if outDir == "" {
		outDir = "grpc/proto"
	}

	conn := values["connection"]
	db := facade.DB(conn)
	for _, table := range strings.Split(values["table"], ",") {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		flag.Infof("开始生成grpc proto: %s", table)
		m.generateProto(db, table, outDir)
	}
}

// generateProto 根据表结构生成grpc proto
func (m *MakeGrpcProto) generateProto(db *gorm.DB, table, outDir string) {
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		flag.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	name := lo.PascalCase(table)
	createFields, userFields, hasCreated, hasUpdated := buildGrpcProtoFields(columns)

	tpl, err := template.ParseFiles(m.GetTemplate("grpc_proto"))
	if err != nil {
		flag.Errorf("解析grpc proto模板失败: %s", err.Error())
		os.Exit(1)
	}

	data := struct {
		Package      string
		GoPackage    string
		Name         string
		ServiceName  string
		CreateFields string
		UserFields   string
		HasCreated   bool
		HasUpdated   bool
	}{
		Package:      "grpc",
		GoPackage:    "gin/" + filepath.ToSlash(outDir) + ";proto",
		Name:         name,
		ServiceName:  name + "Service",
		CreateFields: createFields,
		UserFields:   userFields,
		HasCreated:   hasCreated,
		HasUpdated:   hasUpdated,
	}

	file := filepath.Join(outDir, table+".proto")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	if err = tpl.Execute(f, data); err != nil {
		flag.Errorf("生成grpc proto失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("grpc proto文件: " + file + " 生成成功!")
}

// buildGrpcProtoFields 构建grpc proto字段
func buildGrpcProtoFields(columns []Column) (create, user string, hasCreated, hasUpdated bool) {
	var createLines []string
	var userLines []string

	index := 1
	for _, c := range columns {
		switch c.Name {
		case "id":
			continue
		case model.CreatedField:
			hasCreated = true
			continue
		case model.UpdatedField:
			hasUpdated = true
			continue
		case model.DeletedField:
			continue
		}

		fieldName := grpcProtoFieldName(c.Name)
		fieldType := grpcProtoType(c)
		createLines = append(createLines, fmt.Sprintf("  %s %s = %d;", fieldType, fieldName, index))
		userLines = append(userLines, fmt.Sprintf("  %s %s = %d;", fieldType, fieldName, index+1))
		index++
	}

	if hasCreated {
		userLines = append(userLines, fmt.Sprintf("  string createdAt = %d;", index+1))
		index++
	}
	if hasUpdated {
		userLines = append(userLines, fmt.Sprintf("  string updatedAt = %d;", index+1))
	}

	return strings.Join(createLines, "\n"), strings.Join(userLines, "\n"), hasCreated, hasUpdated
}

// grpcProtoFieldName 数据库字段名转proto字段名
func grpcProtoFieldName(name string) string {
	return lo.CamelCase(lo.PascalCase(name))
}

// grpcProtoType 数据库类型转proto类型
func grpcProtoType(c Column) string {
	t := strings.ToLower(c.DataType)

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
		return "double"
	case t == "json", t == "jsonb":
		return "google.protobuf.Struct"
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		return "string"
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "bytes"
	}

	return "string"
}

func init() {
	cli.Register(&MakeGrpcProto{})
}
