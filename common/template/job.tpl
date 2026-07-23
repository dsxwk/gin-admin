package job

import (
    "gin/app/facade"
    "gin/pkg"
    "gin/pkg/serviceprovider/job"
)

// {{.CamelName}}Job {{.Description}}
type {{.CamelName}}Job struct{}

// {{.CamelName}} {{.Description}}参数
type {{.CamelName}} struct {
}

func (j *{{.CamelName}}Job) Name() string        { return "{{.Name}}" }
func (j *{{.CamelName}}Job) Description() string { return "{{.Description}}" }
func (j *{{.CamelName}}Job) Connection() string  { return "{{.Connection}}" }
func (j *{{.CamelName}}Job) Retry() int          { return {{.Retry}} }
func (j *{{.CamelName}}Job) Delay() int64        { return {{.Delay}} }
func (j *{{.CamelName}}Job) NewPayload() any     { return &{{.CamelName}}{} }

func (j *{{.CamelName}}Job) Handle(payload any) error {
    data := payload.(*{{.CamelName}})
    facade.Log().Info(pkg.Sprintf("Job [{{.Name}}] 处理: %v", data))
    return nil
}

func init() {
    job.Register(&{{.CamelName}}Job{})
}
