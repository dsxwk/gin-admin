package make

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeJob struct {
	base.BaseCommand
}

func (m *MakeJob) Name() string {
	return "make:job"
}

func (m *MakeJob) Description() string {
	return "创建任务Job"
}

func (m *MakeJob) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			Flag:     base.Flag{Short: "n", Long: "name"},
			Desc:     "Job名称, 如 send_email",
			Required: true,
		},
		{
			Flag:     base.Flag{Short: "c", Long: "connection"},
			Desc:     "连接类型: redis, kafka, rabbitmq, sync(默认取queue.connection配置)",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "D", Long: "desc"},
			Desc:     "Job描述",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "R", Long: "retry", Default: "0"},
			Desc:     "重试次数, 默认0(使用全局默认3)",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "d", Long: "delay", Default: "0"},
			Desc:     "重试间隔时间(毫秒), 默认0(使用全局默认1000)",
			Required: false,
		},
	}
}

func (m *MakeJob) Execute(values map[string]string) {
	name := strings.ToLower(values["name"])
	conn := values["connection"]
	desc := values["desc"]

	if conn == "" {
		cfg := facade.Config()
		if cfg != nil {
			conn = cfg.Queue.Connection
		}
		if conn == "" {
			conn = "redis"
		}
	}

	if desc == "" {
		desc = name
	}

	retry, _ := pkg.StringToInt[int](values["retry"])
	delay, _ := pkg.StringToInt[int64](values["delay"])

	camelName := lo.PascalCase(name)

	data := map[string]interface{}{
		"Name":        name,
		"CamelName":   camelName,
		"Connection":  conn,
		"Description": desc,
		"Retry":       retry,
		"Delay":       delay,
	}

	tplFile := filepath.Join(pkg.GetRootPath(), "common", "template", "job.tpl")
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		flag.Errorf("Error parsing job template: %s", err.Error())
		os.Exit(1)
	}

	outFile := filepath.Join("app", "job", name+".go")
	f := m.CheckDirAndFile(outFile)
	if f == nil {
		return
	}

	if err = tpl.Execute(f, data); err != nil {
		flag.Errorf("Error executing job template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("Job文件: %s 生成成功!", outFile)
}

func init() {
	cli.Register(&MakeJob{})
}
