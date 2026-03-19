package make

import (
	"fmt"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"github.com/fatih/color"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MakeQueue struct {
	base.BaseCommand
}

func (m *MakeQueue) Name() string {
	return "make:queue"
}

func (m *MakeQueue) Description() string {
	return "消息队列创建"
}

func (m *MakeQueue) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "t",
				Long:  "type",
			},
			"队列类型, 如: kafka或rabbitmq",
			true,
		},
		{
			base.Flag{
				Short: "n",
				Long:  "name",
			},
			"队列名称, 如: order_create",
			true,
		},
		{
			base.Flag{
				Short: "d",
				Long:  "isDelay",
			},
			"是否延迟队列, 如: true或false",
			false,
		},
		{
			base.Flag{
				Short: "T",
				Long:  "topic",
			},
			"队列主题, 如: kafka_demo",
			false,
		},
		{
			base.Flag{
				Short: "k",
				Long:  "key",
			},
			"消息键, 如: kafka_demo",
			false,
		},
		{
			base.Flag{
				Short: "g",
				Long:  "group",
			},
			"消费组, 如: kafka_demo",
			false,
		},
		{
			base.Flag{
				Short: "q",
				Long:  "queue",
			},
			"队列名, 如: rabbitmq_demo",
			false,
		},
		{
			base.Flag{
				Short: "e",
				Long:  "exchange",
			},
			"交换机, 如: rabbitmq_demo",
			false,
		},
		{
			base.Flag{
				Short: "r",
				Long:  "routing",
			},
			"路由键, 如: rabbitmq_demo",
			false,
		},
		{
			base.Flag{
				Short: "R",
				Long:  "retry",
			},
			"错误重试次数, 如: 3",
			false,
		},
		{
			base.Flag{
				Short: "m",
				Long:  "delayMs",
			},
			"延迟毫秒, 如: 10000",
			false,
		},
	}
}

func (m *MakeQueue) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	color.Green("执行命令: %s %s", m.Name(), m.FormatArgs(values))
	_make := strings.TrimPrefix(m.Name(), "make:")
	if _make != "queue" {
		m.ExitError("命令错误,未找到")
	}
	b := false
	fmt.Printf("values: %v\n", values)
	if values["isDelay"] == "" || values["isDelay"] == "false" {
		b = false
	} else {
		b = true
	}

	m.generateFile(b, values)
}

func init() {
	cli.Register(&MakeQueue{})
}

func (m *MakeQueue) generateFile(isDelay bool, values map[string]string) {
	tmpls := m.GetQueueTemplates()
	tmpl, err := template.ParseFiles(tmpls["consumer"], tmpls["producer"])
	if err != nil {
		color.Red("Error parsing template:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Defined templates:", tmpl.DefinedTemplates())

	files := m.GetMakeQueueFile(values["name"], values["type"], isDelay)
	file1 := files["consumer"]
	file2 := files["producer"]

	// 创建文件
	f1 := m.CheckDirAndFile(file1)
	if f1 == nil {
		os.Exit(1)
	}

	f2 := m.CheckDirAndFile(file2)
	if f2 == nil {
		os.Exit(1)
	}

	generateTemplateData(isDelay, values, tmpls, f1, f2, tmpl)

	color.Green(flag.Success + "  验证请求文件: " + file1 + " 生成成功!")
	color.Green(flag.Success + "  验证请求文件: " + file2 + " 生成成功!")
}

func generateTemplateData(isDelay bool, values, tmpls map[string]string, f1, f2 io.Writer, tmpl *template.Template) (any, any) {
	// 包名
	packageName1 := "consumer"
	packageName2 := "producer"

	retry, err := strconv.Atoi(values["retry"])
	if err != nil {
		retry = 3
	}

	delayMs, err := strconv.ParseInt(values["delayMs"], 10, 64)
	if err != nil {
		delayMs = 0
	}

	data1 := struct {
		Package  string // 提取的包名
		Name     string // 模块名称(首字母大写)
		Type     string // 队列类型
		Topic    string // 队列主题
		Group    string // 消费组
		Queue    string // 队列名称
		Exchange string // 交换机
		Routing  string // 路由键
		Retry    int    // 重试次数
		IsDelay  bool   // 是否延迟队列

	}{
		Package:  packageName1,
		Name:     pkg.ToUpperCamel(values["name"]),
		Type:     values["type"],
		Topic:    values["topic"],
		Group:    values["group"],
		Queue:    values["queue"],
		Exchange: values["exchange"],
		Routing:  values["routing"],
		Retry:    retry,
		IsDelay:  isDelay,
	}

	err = tmpl.ExecuteTemplate(f1, filepath.Base(tmpls["consumer"]), data1)
	if err != nil {
		color.Red("Error executing consumer template:", err.Error())
		os.Exit(1)
	}

	data2 := struct {
		Package  string // 提取的包名
		Name     string // 模块名称(首字母大写)
		Type     string // 队列类型
		Topic    string // 队列主题
		Key      string // 消息键
		Queue    string // 队列名称
		Exchange string // 交换机
		Routing  string // 路由键
		IsDelay  bool   // 是否延迟队列
		DelayMs  int64  // 延迟毫秒
	}{
		Package:  packageName2,
		Name:     pkg.ToUpperCamel(values["name"]),
		Type:     values["type"],
		Topic:    values["topic"],
		Key:      values["group"],
		Queue:    values["queue"],
		Exchange: values["exchange"],
		Routing:  values["routing"],
		IsDelay:  isDelay,
		DelayMs:  delayMs,
	}

	err = tmpl.ExecuteTemplate(f2, filepath.Base(tmpls["producer"]), data2)
	if err != nil {
		color.Red("Error executing producer template:", err.Error())
		os.Exit(1)
	}

	return data1, data2
}
