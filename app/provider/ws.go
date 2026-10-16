package provider

import (
	"context"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/ws"
	"gin/websocket"
	"time"
)

// WsProvider WebSocket服务提供者
type WsProvider struct {
	manager   *ws.Manager
	closeWait time.Duration
}

// Name 服务提供者名称
func (p *WsProvider) Name() string {
	return serviceprovider.ServiceWS
}

// Register 注册WebSocket服务
func (p *WsProvider) Register(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Ws.Enabled {
		return
	}

	p.manager = ws.NewManager(ws.Options{
		Shards:          cfg.Ws.Shards,
		MaxConnections:  cfg.Ws.MaxConnections,
		WriteQueueSize:  cfg.Ws.WriteQueueSize,
		ReadBufferSize:  cfg.Ws.ReadBufferSize,
		WriteBufferSize: cfg.Ws.WriteBufferSize,
		ReadLimit:       cfg.Ws.ReadLimit,
		WriteWait:       cfg.Ws.WriteWait,
		PongWait:        cfg.Ws.PongWait,
		PingInterval:    cfg.Ws.PingInterval,
		CloseWait:       cfg.Ws.CloseWait,
		AllowedOrigins:  cfg.Ws.AllowedOrigins,
		Compression:     cfg.Ws.Compression,
	}, websocket.All()...)
	p.closeWait = cfg.Ws.CloseWait
	if p.closeWait <= 0 {
		p.closeWait = 5 * time.Second
	}

	app.SetWS(p.manager)
	path := cfg.Ws.Path
	if path == "" {
		path = "/ws"
	}
	flag.Infof("ws服务注册成功,地址: %s", path)
}

// Boot 启动WebSocket服务
func (p *WsProvider) Boot(_ *container.Container) {
	// WebSocket由HTTP路由触发升级
}

// Runners 后台运行任务
func (p *WsProvider) Runners() []serviceprovider.Runner {
	if p.manager == nil {
		return nil
	}

	manager := p.manager
	closeWait := p.closeWait
	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("websocket_shutdown", func() error {
			ctx, cancel := context.WithTimeout(context.Background(), closeWait)
			defer cancel()
			return manager.Close(ctx)
		}),
	}
}

// Dependencies 依赖服务
func (p *WsProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
