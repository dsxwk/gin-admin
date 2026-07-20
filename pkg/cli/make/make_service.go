package make

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeService struct {
	base.BaseCommand
}

func (m *MakeService) Name() string {
	return "make:service"
}

func (m *MakeService) Description() string {
	return "服务创建"
}

func (m *MakeService) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: v1/user",
			true,
		},
		{
			base.Flag{
				Short: "t",
				Long:  "table",
			},
			"表名, 用于生成模型字段",
			false,
		},
		{
			base.Flag{
				Short:   "c",
				Long:    "connection",
				Default: "mysql",
			},
			"数据库连接",
			false,
		},
	}
}

func (m *MakeService) Execute(values map[string]string) {
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)

	conn := values["connection"]
	table := values["table"]

	var fields []Fields
	// 获取表字段
	if table != "" {
		db := facade.DB(conn)
		columns, err := GetColumnInfo(db, table)
		if err != nil {
			facade.Log().Warn(err.Error())
		} else {
			for _, col := range columns {
				// 排除字段
				if col.Name == "id" ||
					col.Name == "created_at" ||
					col.Name == "updated_at" ||
					col.Name == "deleted_at" {
					continue
				}
				name := m.toGoName(col.Name)
				fields = append(fields, Fields{
					Name:   name,
					IsJSON: m.isJSONType(col.ColumnType),
				})
			}
		}
	}

	m.generateFile(_make, f, fields)
}

func init() {
	cli.Register(&MakeService{})
}

// Fields 字段信息
type Fields struct {
	Name   string // Go字段名
	IsJSON bool   // 是否是JSON字段
}

// isJSONType 判断是否是JSON类型
func (m *MakeService) isJSONType(dbType string) bool {
	t := strings.ToLower(dbType)
	return t == "json" || t == "jsonb"
}

// toGoName 转Go命名
func (m *MakeService) toGoName(name string) string {
	if name == "id" {
		return "ID"
	}
	return lo.PascalCase(name)
}

func (m *MakeService) generateFile(_make, file string, fields []Fields) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}

	// 提取包名 (文件路径中的最后一个目录作为包名)
	packageName := filepath.Base(filepath.Dir(file))

	// 创建文件
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	data := struct {
		Package   string // 提取的包名
		Name      string // 模块名称(首字母大写)
		Fields    string
		HasFields bool
	}{
		Package:   packageName,
		Name:      lo.PascalCase(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Fields:    m.buildModelFields(fields),
		HasFields: fields != nil && len(fields) > 0,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("服务文件: " + file + " 生成成功!")
}

// buildModelFields 构建模型字段
func (m *MakeService) buildModelFields(fields []Fields) string {
	if fields == nil || len(fields) == 0 {
		return ""
	}

	// 计算所有字段名的最大长度
	maxNameLen := 0
	for _, f := range fields {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
	}

	var lines []string
	for _, f := range fields {
		// 计算需要填充的空格数:最大长度-当前字段名长度
		padding := strings.Repeat(" ", maxNameLen-len(f.Name))
		var line string
		if f.IsJSON {
			line = fmt.Sprintf("\t\t%s:%s &model.JsonValue{Data: req.%s},", f.Name, padding, f.Name)
		} else {
			line = fmt.Sprintf("\t\t%s:%s req.%s,", f.Name, padding, f.Name)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
