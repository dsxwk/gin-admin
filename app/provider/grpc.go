package provider

import (
	"context"
	"fmt"
	"gin/common/flag"
	grpcservice "gin/grpc/service"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/grpcclient"
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
func (p *GrpcProvider) Register(app *container.Container) {
}

// Boot 启动服务
func (p *GrpcProvider) Boot(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Grpc.Enabled {
		return
	}

	client, err := grpcclient.NewClient(fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		flag.Errorf("grpc客户端创建失败: %v", err)
		return
	}

	srv, err := grpcclient.NewServer(
		cfg.Grpc.Host,
		cfg.Grpc.Port,
		cfg.Jwt.Key,
		grpcservice.Services()...,
	)
	if err != nil {
		_ = client.Close()
		flag.Errorf("grpc服务启动失败: %v", err)
		return
	}

	if err = srv.Start(); err != nil {
		_ = srv.Stop()
		_ = client.Close()
		flag.Errorf("grpc服务启动失败: %v", err)
		return
	}

	p.client = client
	p.server = srv
	app.SetGRPC(client)
	flag.Infof("grpc服务启动成功: %s", srv.Addr())
}

// Runners 后台运行任务
func (p *GrpcProvider) Runners() []serviceprovider.Runner {
	if p.server == nil && p.client == nil {
		return nil
	}
	return []serviceprovider.Runner{
		&grpcShutdownRunner{
			server: p.server,
			client: p.client,
		},
	}
}

// Dependencies 依赖服务
func (p *GrpcProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog, serviceprovider.ServiceDB}
}

// grpcShutdownRunner grpc关闭任务
type grpcShutdownRunner struct {
	server *grpcclient.Server
	client *grpcclient.Client
}

// Run 运行等待任务
func (r *grpcShutdownRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 停止服务
func (r *grpcShutdownRunner) Stop() error {
	var serverErr error
	if r.server != nil {
		serverErr = r.server.Stop()
	}
	var clientErr error
	if r.client != nil {
		clientErr = r.client.Close()
	}
	if serverErr != nil {
		return serverErr
	}
	return clientErr
}

// Name 任务名称
func (r *grpcShutdownRunner) Name() string {
	return "grpc_shutdown"
}
