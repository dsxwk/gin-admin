package provider

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/filesystem"
	"gin/pkg/terminal"
)

// FileProvider 文件系统服务提供者
type FileProvider struct {
	manager *filesystem.Manager
}

// Name 服务提供者名称
func (p *FileProvider) Name() string {
	return serviceprovider.ServiceFile
}

// Register 注册服务到容器
func (p *FileProvider) Register(app *container.Container) {
	manager, err := filesystem.NewManager(app.Config().Filesystem)
	if err != nil {
		terminal.Errorf("文件系统服务注册失败: %v", err)
		return
	}

	p.manager = manager
	app.SetFile(manager)
}

// Boot 启动服务
func (p *FileProvider) Boot(_ *container.Container) {
	if p.manager == nil {
		return
	}
	terminal.Infof("文件系统服务启动成功")
}

// Dependencies 依赖服务
func (p *FileProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
