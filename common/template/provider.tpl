package {{.Package}}

import (
{{- if .HasRunner}}
	"context"
{{- end}}
	"gin/app/facade"
	"gin/pkg/foundation"
)

func init() {
	foundation.Register(&{{.ProviderName}}Provider{})
}

// {{.ProviderName}}Provider {{.Desc}}服务提供者
type {{.ProviderName}}Provider struct{}

// Name 服务提供者名称
func (p *{{.ProviderName}}Provider) Name() string {
	return "{{.ProviderVar}}"
}

// Register 注册服务到门面
func (p *{{.ProviderName}}Provider) Register(app foundation.App) {
	// TODO: 注册服务到门面
	// facade.Register("{{.ProviderVar}}", facade.{{.ProviderName}}Provider)
}

// Boot 启动服务
func (p *{{.ProviderName}}Provider) Boot(app foundation.App) {
	// TODO: 初始化服务
	facade.Log.Info("{{.ProviderName}}服务启动成功")
}
{{- if .HasRunner}}

// Runners 后台运行任务(用于优雅关闭)
func (p *{{.ProviderName}}Provider) Runners() []foundation.Runner {
	return []foundation.Runner{
		&{{.ProviderName}}CleanupRunner{},
	}
}
{{- end}}

// Dependencies 依赖服务
func (p *{{.ProviderName}}Provider) Dependencies() []string {
	return []string{ {{- range $i, $dep := .Deps}}{{if $i}}, {{end}}"{{$dep}}"{{- end}} }
}
{{- if .HasRunner}}

// {{.ProviderName}}CleanupRunner {{.ProviderVar}}清理任务
type {{.ProviderName}}CleanupRunner struct{}

// Run 运行清理任务
func (r *{{.ProviderName}}CleanupRunner) Run(ctx context.Context) error {
	// 等待停止信号
    <-ctx.Done()
    return nil
}

// Stop 停止时清理资源
func (r *{{.ProviderName}}CleanupRunner) Stop() error {
	// TODO: 清理资源
	facade.Log.Info("{{.ProviderVar}}服务已关闭")
	return nil
}

// Name 任务名称
func (r *{{.ProviderName}}CleanupRunner) Name() string {
	return "{{.ProviderVar}}_cleanup"
}
{{- end}}