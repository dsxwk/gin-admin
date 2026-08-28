package make

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// MakeGrpcModel grpc模型生成命令
type MakeGrpcModel struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeGrpcModel) Name() string {
	return "grpc-make:model"
}

// Description 命令描述
func (m *MakeGrpcModel) Description() string {
	return "grpc模型创建"
}

// Help 命令选项
func (m *MakeGrpcModel) Help() []base.CommandOption {
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
				Default: "grpc/model",
			},
			"输出目录, 默认grpc/model",
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
func (m *MakeGrpcModel) Execute(values map[string]string) {
	outDir := strings.TrimPrefix(values["path"], "/")
	if outDir == "" {
		outDir = "grpc/model"
	}

	conn := values["connection"]
	db := facade.DB(conn)
	for _, table := range strings.Split(values["table"], ",") {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		flag.Infof("开始生成grpc模型: %s", table)
		m.generateModel(db, table, outDir)
	}
}

// generateModel 根据表结构生成Go模型
func (m *MakeGrpcModel) generateModel(db *gorm.DB, table, outDir string) {
	columns, err := GetColumnInfo(db, table)
	if err != nil {
		flag.Errorf("获取表字段失败: %s", err.Error())
		os.Exit(1)
	}

	tableComment, _ := getTableComment(db, table)
	imports := make(map[string]string)
	structName := lo.PascalCase(table)
	tableConst := "TableName" + structName

	maxNameLen := 0
	maxTypeLen := 0
	for _, c := range columns {
		fieldName := grpcModelFieldName(c.Name)
		fieldType := grpcGoType(c, imports)
		if len(fieldName) > maxNameLen {
			maxNameLen = len(fieldName)
		}
		if len(fieldType) > maxTypeLen {
			maxTypeLen = len(fieldType)
		}
	}

	var fieldLines []string
	for _, c := range columns {
		fieldName := grpcModelFieldName(c.Name)
		fieldType := grpcGoType(c, imports)
		jsonName := lo.CamelCase(fieldName)
		gormTag := buildGormTag(c)
		tag := "`" + gormTag + " json:\"" + jsonName + "\" form:\"" + jsonName + "\"`"
		if jsonName == "deletedAt" || c.Name == "deleted_at" {
			tag = strings.TrimSuffix(tag, "`") + " swaggerignore:\"true\"`"
		}
		line := fmt.Sprintf("%s%s%s%s%s",
			fieldName,
			strings.Repeat(" ", maxNameLen-len(fieldName)+1),
			fieldType,
			strings.Repeat(" ", maxTypeLen-len(fieldType)+1),
			tag,
		)
		fieldLines = append(fieldLines, strings.TrimSpace(line))
	}

	tpl, err := template.ParseFiles(m.GetTemplate("model"))
	if err != nil {
		flag.Errorf("解析模型模板失败: %s", err.Error())
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
		Imports:       renderGrpcImports(imports),
		Struct:        structName,
		StructComment: tableComment,
		Table:         table,
		TableConst:    tableConst,
		Connection:    "mysql",
		Fields:        fieldLines,
	}

	file := filepath.Join(outDir, table+".go")
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}
	if err = tpl.Execute(f, data); err != nil {
		flag.Errorf("生成grpc模型失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("grpc模型文件: " + file + " 生成成功!")
}

// grpcGoType 数据库类型转grpc模型Go类型
func grpcGoType(c Column, imports map[string]string) string {
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
		return "float64"
	case t == "json", t == "jsonb":
		return "*JsonValue"
	case strings.Contains(t, "timestamp"),
		strings.Contains(t, "datetime"),
		t == "date":
		if c.Name == "deleted_at" {
			return "*DeletedAt"
		}
		return "*DateTime"
	case strings.Contains(t, "blob"),
		strings.Contains(t, "binary"),
		strings.Contains(t, "bytea"):
		return "[]byte"
	}

	return "string"
}

// renderGrpcImports 渲染import块
func renderGrpcImports(imports map[string]string) string {
	if len(imports) == 0 {
		return ""
	}

	keys := make([]string, 0, len(imports))
	for pkg := range imports {
		keys = append(keys, pkg)
	}
	sort.Strings(keys)

	lines := []string{"import ("}
	for _, pkg := range keys {
		if alias := imports[pkg]; alias != "" {
			lines = append(lines, fmt.Sprintf("\t%s %q", alias, pkg))
		} else {
			lines = append(lines, fmt.Sprintf("\t%q", pkg))
		}
	}
	lines = append(lines, ")")
	return strings.Join(lines, "\n")
}

// grpcModelFieldName 数据库字段名转Go字段名
func grpcModelFieldName(name string) string {
	if name == "id" {
		return "ID"
	}
	return lo.PascalCase(name)
}

func init() {
	cli.Register(&MakeGrpcModel{})
}
