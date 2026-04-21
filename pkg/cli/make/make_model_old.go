package make

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MakeModelOld struct {
	base.BaseCommand
}

func (m *MakeModelOld) Name() string {
	return "make:model-old"
}

func (m *MakeModelOld) Description() string {
	return "模型创建old"
}

func (m *MakeModelOld) Help() []base.CommandOption {
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
			"是否驼峰字段, 如: true",
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

func (m *MakeModelOld) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	// 去除前斜杠
	p := filepath.Join("app/model/", strings.TrimPrefix(values["path"], "/"))
	tables := strings.Split(values["table"], ",")
	for i := range tables {
		tables[i] = strings.TrimSpace(tables[i])
		flag.Successf("创建模型: %s (表名: %s 是否使用驼峰: %v)\n", p+"/"+tables[i]+".gen.go", tables[i], values["camel"])
	}

	m.generateFiles(p, values["connection"], tables, m.StringToBool(values["camel"]))
}

func init() {
	cli.Register(&MakeModelOld{})
}

// generateFiles 生成模型文件
func (m *MakeModelOld) generateFiles(path, conn string, tables []string, camel bool) {
	var (
		root    = pkg.GetRootPath()
		p       = filepath.Base(path)
		outPath = filepath.Join(root + "/app/temp")
	)

	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
		ModelPkgPath:      path,
	})

	g.UseDB(facade.DB.Connection())

	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"json": func(detailType gorm.ColumnType) (dataType string) {
			if p != "model" {
				return "*model.JsonValue"
			} else {
				return "*JsonValue"
			}
		},
		"datetime": func(detailType gorm.ColumnType) (dataType string) {
			// deleted_at字段特殊处理
			if detailType.Name() == "deleted_at" {
				if p != "model" {
					return "*model.DeletedAt"
				} else {
					return "*DeletedAt"
				}
			}

			if p != "model" {
				return "*model.DateTime"
			} else {
				return "*DateTime"
			}
		},
		// "timestamp":  func(detailType gorm.ColumnType) (dataType string) { return "string" }, // 添加此行将 timestamp 转换为 string
		// "date":       func(detailType gorm.ColumnType) (dataType string) { return "string" }, // 添加此行将 date 转换为 string
	}

	// 要先于`ApplyBasic`执行
	g.WithDataTypeMap(dataMap)

	// 自定义JSON tag
	if camel {
		g.WithJSONTagNameStrategy(func(columnName string) string {
			return lo.CamelCase(columnName)
		})
	}

	flag.Infof("开始生成模型, 共 %d 张表", len(tables))

	for _, table := range tables {
		flag.Infof("→ 正在生成表: %s", table)
		fileName := filepath.Join(path, table+".gen.go")
		f := m.CheckDirAndFile(fileName)
		if f == nil {
			continue
		}
		model := g.GenerateModel(table)
		g.ApplyBasic(model)
		g.Execute()
		// 自动追加swaggerignore:"true"
		content, err := os.ReadFile(fileName)
		if err != nil {
			continue
		}
		text := string(content)

		re := regexp.MustCompile("(`[^`]*json:\"deletedAt\"[^`]*`)")

		text = re.ReplaceAllStringFunc(text, func(match string) string {
			if strings.Contains(match, "swaggerignore") {
				return match
			}
			return strings.TrimSuffix(match, "`") + " swaggerignore:\"true\"`"
		})

		if err = os.WriteFile(fileName, []byte(text), 0644); err != nil {
			flag.Errorf("为文件 %s 添加 swaggerignore 失败", fileName)
			os.Exit(1)
		}

		// 为每个生成的模型文件追加Connection方法
		if conn != "" {
			structName := lo.CamelCase(table)
			err = appendConnection(path, table, structName, conn)
			if err != nil {
				flag.Errorf("为模型 %s 追加 Connection 方法失败: %s", table, err.Error())
			}
		}

		flag.Successf("模型文件: " + fileName + " 生成成功!")
	}

	_ = os.RemoveAll(outPath)
}

// appendConnection 为模型文件追加Connection方法
func appendConnection(path, table, structName, conn string) error {
	filePath := filepath.Join(path, table+".gen.go")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(content)
	// 检查是否已经存在Connection方法
	if strings.Contains(text, "func (*"+structName+") Connection() string") {
		flag.Infof("模型 %s 已存在 Connection 方法,跳过添加", structName)
		return nil
	}

	// 在TableName方法后面追加Connection方法
	function := fmt.Sprintf(`
// Connection 数据库连接名称
func (*%s) Connection() string {
    return "%s"
}`, structName, conn)

	// 查找TableName方法结束的位置
	tableNamePattern := regexp.MustCompile(`func \(\*` + structName + `\) TableName\(\) string \{[\s\S]*?\n\}`)

	if tableNamePattern.MatchString(text) {
		// 在TableName方法后面插入
		text = tableNamePattern.ReplaceAllStringFunc(text, func(match string) string {
			return match + "\n" + function
		})
	}

	return os.WriteFile(filePath, []byte(text), 0644)
}
