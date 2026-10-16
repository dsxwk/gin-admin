package job

import (
	"context"
	"gin/app/facade"
	"gin/pkg"
	servicejob "gin/pkg/serviceprovider/job"
)

// {{.CamelName}} {{.Description}}参数
type {{.CamelName}} struct {
}

// New{{.CamelName}}Job 创建{{.Description}}
func New{{.CamelName}}Job() *servicejob.Adapter[{{.CamelName}}] {
	return servicejob.NewAdapter(
		servicejob.Config{
			Name:        "{{.Name}}",
			Description: "{{.Description}}",
			Connection:  "{{.Connection}}",
			Retry:       {{.Retry}},
			Delay:       {{.Delay}},
		},
		handle{{.CamelName}},
	)
}

// handle{{.CamelName}} 处理{{.Description}}
func handle{{.CamelName}}(_ context.Context, data *{{.CamelName}}) error {
	facade.Log().Info(pkg.Sprintf("Job [{{.Name}}] 处理: %v", data))
	// 待处理业务逻辑
	return nil
}
