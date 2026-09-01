package provider

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/common/flag"
	_ "gin/grpc/service"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/grpcclient"
)

func init() {
	serviceprovider.Register(&GrpcProvider{})
}

// GrpcProvider grpc服务提供者
type GrpcProvider struct {
	server *grpcclient.Server
	client *grpcclient.Client
}

// Name 服务提供者名称
func (p *GrpcProvider) Name() string {
	return "grpc"
}

// Register 注册服务到门面
func (p *GrpcProvider) Register(app serviceprovider.App) {
	cfg := facade.Config()
	if cfg == nil || !cfg.Grpc.Enabled {
		return
	}

	c, err := grpcclient.NewClient(fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		flag.Errorf("grpc客户端创建失败: %v", err)
		return
	}
	p.client = c
	facade.Register[*grpcclient.Client]("grpc", c)
}

// Boot 启动服务
func (p *GrpcProvider) Boot(app serviceprovider.App) {
	cfg := facade.Config()
	if cfg == nil || !cfg.Grpc.Enabled {
		return
	}

	grpcclient.SetJwtKey(cfg.Jwt.Key)
	srv, err := grpcclient.NewServer(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		flag.Errorf("grpc服务启动失败: %v", err)
		return
	}

	p.server = srv
	if err = srv.Start(); err != nil {
		flag.Errorf("grpc服务启动失败: %v", err)
		return
	}
	flag.Infof("grpc服务启动成功: %s", srv.Addr())
}

// Runners 后台运行任务
func (p *GrpcProvider) Runners() []serviceprovider.Runner {
	if p.server == nil {
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
	return []string{"config", "log", "db"}
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
	if r.server != nil {
		if err := r.server.Stop(); err != nil {
			return err
		}
	}
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Name 任务名称
func (r *grpcShutdownRunner) Name() string {
	return "grpc_shutdown"
}
