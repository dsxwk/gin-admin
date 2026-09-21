package container

import (
	"gin/config"
	"gin/pkg/serviceprovider/cache"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/es"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/serviceprovider/grpcclient"
	"gin/pkg/serviceprovider/http"
	"gin/pkg/serviceprovider/job"
	"gin/pkg/serviceprovider/lang"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/mcp"
	"gin/pkg/serviceprovider/orm"
	"gin/pkg/serviceprovider/queue"
	"gin/pkg/serviceprovider/ratelimit"
	"gin/pkg/serviceprovider/request"
	"sync/atomic"
)

var app = New()

// Default 获取默认服务容器
func Default() *Container {
	return app
}

// Container 服务容器
type Container struct {
	config    atomic.Pointer[config.Config]
	log       atomic.Pointer[logger.Logger]
	db        atomic.Pointer[orm.Manager]
	cache     atomic.Pointer[cache.Manager]
	event     atomic.Pointer[eventbus.Registry]
	debugger  atomic.Pointer[debugger.Debugger]
	http      atomic.Pointer[http.Client]
	es        atomic.Pointer[es.Client]
	grpc      atomic.Pointer[grpcclient.Client]
	queue     atomic.Pointer[queue.Manager]
	job       atomic.Pointer[job.Manager]
	lang      atomic.Pointer[lang.Service]
	mcp       atomic.Pointer[mcp.Handler]
	request   atomic.Pointer[request.Client]
	rateLimit atomic.Pointer[ratelimit.Manager]
}

// New 创建服务容器
func New() *Container {
	return &Container{}
}

// SetConfig 设置配置服务
func (c *Container) SetConfig(value *config.Config) {
	c.config.Store(value)
}

// Config 获取配置服务
func (c *Container) Config() *config.Config {
	return c.config.Load()
}

// SetLog 设置日志服务
func (c *Container) SetLog(value *logger.Logger) {
	c.log.Store(value)
}

// Log 获取日志服务
func (c *Container) Log() *logger.Logger {
	return c.log.Load()
}

// SetDB 设置数据库服务
func (c *Container) SetDB(value *orm.Manager) {
	c.db.Store(value)
}

// DB 获取数据库服务
func (c *Container) DB() *orm.Manager {
	return c.db.Load()
}

// SetCache 设置缓存服务
func (c *Container) SetCache(value *cache.Manager) {
	c.cache.Store(value)
}

// Cache 获取缓存服务
func (c *Container) Cache() *cache.Manager {
	return c.cache.Load()
}

// Redis 获取Redis服务
func (c *Container) Redis() *cache.RedisCache {
	manager := c.Cache()
	if manager == nil {
		return nil
	}
	return manager.Redis()
}

// SetEvent 设置事件服务
func (c *Container) SetEvent(value *eventbus.Registry) {
	c.event.Store(value)
}

// Event 获取事件服务
func (c *Container) Event() *eventbus.Registry {
	return c.event.Load()
}

// SetDebugger 设置调试服务
func (c *Container) SetDebugger(value *debugger.Debugger) {
	c.debugger.Store(value)
}

// Debugger 获取调试服务
func (c *Container) Debugger() *debugger.Debugger {
	return c.debugger.Load()
}

// SetHTTP 设置HTTP服务
func (c *Container) SetHTTP(value *http.Client) {
	c.http.Store(value)
}

// HTTP 获取HTTP服务
func (c *Container) HTTP() *http.Client {
	return c.http.Load()
}

// SetES 设置ES服务
func (c *Container) SetES(value *es.Client) {
	c.es.Store(value)
}

// ES 获取ES服务
func (c *Container) ES() *es.Client {
	return c.es.Load()
}

// SetGRPC 设置GRPC服务
func (c *Container) SetGRPC(value *grpcclient.Client) {
	c.grpc.Store(value)
}

// GRPC 获取GRPC服务
func (c *Container) GRPC() *grpcclient.Client {
	return c.grpc.Load()
}

// SetQueue 设置队列服务
func (c *Container) SetQueue(value *queue.Manager) {
	c.queue.Store(value)
}

// Queue 获取队列服务
func (c *Container) Queue() *queue.Manager {
	return c.queue.Load()
}

// SetJob 设置任务服务
func (c *Container) SetJob(value *job.Manager) {
	c.job.Store(value)
}

// Job 获取任务服务
func (c *Container) Job() *job.Manager {
	return c.job.Load()
}

// SetLang 设置翻译服务
func (c *Container) SetLang(value *lang.Service) {
	c.lang.Store(value)
}

// Lang 获取翻译服务
func (c *Container) Lang() *lang.Service {
	return c.lang.Load()
}

// SetMCP 设置MCP服务
func (c *Container) SetMCP(value *mcp.Handler) {
	c.mcp.Store(value)
}

// MCP 获取MCP服务
func (c *Container) MCP() *mcp.Handler {
	return c.mcp.Load()
}

// SetRequest 设置请求服务
func (c *Container) SetRequest(value *request.Client) {
	c.request.Store(value)
}

// Request 获取请求服务
func (c *Container) Request() *request.Client {
	return c.request.Load()
}

// SetRateLimit 设置限流服务
func (c *Container) SetRateLimit(value *ratelimit.Manager) {
	c.rateLimit.Store(value)
}

// RateLimit 获取限流服务
func (c *Container) RateLimit() *ratelimit.Manager {
	return c.rateLimit.Load()
}
