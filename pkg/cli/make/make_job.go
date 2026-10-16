package make

import (
	"fmt"
	"gin/app/facade"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"path/filepath"
	"strings"
	"text/template"

	"gin/pkg/convert"
	"gin/pkg/pathx"
	"github.com/samber/lo"
)

type MakeJob struct {
	base.Command
}

// NewMakeJob 创建命令
func NewMakeJob() *MakeJob {
	return &MakeJob{
		Command: base.NewCommand(base.Config{
			Name:        "make:job",
			Description: "创建任务Job",
			Options: []base.CommandOption{
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
					Desc:     "重试次数, 0表示不重试",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "delay", Default: "0"},
					Desc:     "首次延迟及重试间隔(毫秒), 0表示不延迟",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeJob) Handle(values map[string]string) error {
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

	retry, _ := convert.StringToInt[int](values["retry"])
	delay, _ := convert.StringToInt[int64](values["delay"])

	camelName := lo.PascalCase(name)

	data := map[string]any{
		"Name":        name,
		"CamelName":   camelName,
		"Connection":  conn,
		"Description": desc,
		"Retry":       retry,
		"Delay":       delay,
	}

	tplFile := filepath.Join(pathx.RootPath(), "pkg", "cli", "template", "job.tpl")
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		terminal.Errorf("Error parsing job template: %s", err.Error())
		return fmt.Errorf("命令执行失败")

	}

	outFile := filepath.Join("app", "job", name+".go")
	f := m.CheckDirAndFile(outFile)
	if f == nil {
		return nil

	}

	if err = tpl.Execute(f, data); err != nil {
		terminal.Errorf("Error executing job template: %s", err.Error())
		return fmt.Errorf("命令执行失败")

	}

	jobsFile := filepath.Join("app", "job", "registry.go")
	if err = cli.AddItem(jobsFile, "return []servicejob.Job{", "New"+camelName+"Job()"); err != nil {
		terminal.Errorf("自动注册Job失败: %s", err.Error())
		return nil

	}

	terminal.Successf("Job文件: %s 生成成功!", outFile)
	return nil
}
