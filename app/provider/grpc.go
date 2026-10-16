package provider

import (
	"fmt"
	grpcservice "gin/grpc/service"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/grpcclient"
	"gin/pkg/terminal"
)

// GrpcProvider grpc服务提供者
type GrpcProvider struct {
	server *grpcclient.Server
	client *grpcclient.Client
}

// Name 服务提供者名称
func (p *GrpcProvider) Name() string {
	return serviceprovider.ServiceGRPC
}

// Register 注册服务到容器
func (p *GrpcProvider) Register(_ *container.Container) {
}

// Boot 启动服务
func (p *GrpcProvider) Boot(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Grpc.Enabled {
		return
	}

	client, err := grpcclient.NewClient(fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		terminal.Errorf("grpc客户端创建失败: %v", err)
		return
	}

	srv, err := grpcclient.NewServer(
		cfg.Grpc.Host,
		cfg.Grpc.Port,
		cfg.Jwt.Key,
		grpcservice.All()...,
	)
	if err != nil {
		_ = client.Close()
		terminal.Errorf("grpc服务启动失败: %v", err)
		return
	}

	if err = srv.Start(); err != nil {
		_ = srv.Stop()
		_ = client.Close()
		terminal.Errorf("grpc服务启动失败: %v", err)
		return
	}

	p.client = client
	p.server = srv
	app.SetGRPC(client)
	terminal.Infof("grpc服务启动成功: %s", srv.Addr())
}

// Runners 后台运行任务
func (p *GrpcProvider) Runners() []serviceprovider.Runner {
	if p.server == nil && p.client == nil {
		return nil
	}

	server := p.server
	client := p.client
	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("grpc_shutdown", func() error {
			var serverErr error
			if server != nil {
				serverErr = server.Stop()
			}
			var clientErr error
			if client != nil {
				clientErr = client.Close()
			}
			if serverErr != nil {
				return serverErr
			}
			return clientErr
		}),
	}
}

// Dependencies 依赖服务
func (p *GrpcProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog, serviceprovider.ServiceDB}
}
