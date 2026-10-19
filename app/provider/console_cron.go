package provider

import (
	"gin/app/console"
	"gin/pkg/cli/base"
	"gin/pkg/cli/registry"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/terminal"
)

// ConsoleCronProvider 控制台定时任务提供者
type ConsoleCronProvider struct {
	manager *base.CronManager
}

// Name 服务提供者名称
func (p *ConsoleCronProvider) Name() string {
	return serviceprovider.ServiceConsoleCron
}

// Register 注册控制台定时任务
func (p *ConsoleCronProvider) Register(_ *container.Container) {
	commands := append(console.All(), registry.All()...)
	manager, err := base.NewCronManager(commands...)
	if err != nil {
		terminal.Errorf("控制台定时任务注册失败: %s", err.Error())
		return
	}
	p.manager = manager
}

// Boot 启动控制台定时任务
func (p *ConsoleCronProvider) Boot(_ *container.Container) {
	if p.manager == nil || p.manager.Len() == 0 {
		return
	}
	terminal.Infof("控制台定时任务注册成功,任务数: %d", p.manager.Len())
}

// Runners 后台运行任务
func (p *ConsoleCronProvider) Runners() []serviceprovider.Runner {
	if p.manager == nil || p.manager.Len() == 0 {
		return nil
	}
	return []serviceprovider.Runner{p.manager}
}

// Dependencies 依赖服务
func (p *ConsoleCronProvider) Dependencies() []string {
	return []string{
		serviceprovider.ServiceConfig,
		serviceprovider.ServiceLog,
		serviceprovider.ServiceAgent,
		serviceprovider.ServiceDB,
		serviceprovider.ServiceCache,
		serviceprovider.ServiceFile,
		serviceprovider.ServiceEvent,
		serviceprovider.ServiceDebugger,
		serviceprovider.ServiceHTTP,
		serviceprovider.ServiceES,
		serviceprovider.ServiceGRPC,
		serviceprovider.ServiceQueue,
		serviceprovider.ServiceJob,
		serviceprovider.ServiceLang,
		serviceprovider.ServiceMCP,
		serviceprovider.ServiceWS,
		serviceprovider.ServiceRequest,
		serviceprovider.ServiceRateLimit,
	}
}
