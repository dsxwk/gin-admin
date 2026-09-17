package job

import (
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	"os"
)

// Job 任务接口
type Job interface {
	queue.Named
	queue.PayloadHandler
	// Description 任务描述
	Description() string
	// Connection 任务连接, 默认 "redis",可选 "sync","kafka","rabbitmq"
	Connection() string
	// Retry 重试次数, 默认3
	Retry() int
	// Delay 重试间隔时间(毫秒), 默认1000
	Delay() int64
}

// Registry 任务注册表
type Registry struct {
	items *queue.Registry[Job]
}

// NewRegistry 创建任务注册表
func NewRegistry() *Registry {
	return &Registry{
		items: queue.NewRegistry[Job](),
	}
}

var registry = NewRegistry()

// Register 注册任务
func (r *Registry) Register(item Job) error {
	return r.items.Register(item)
}

// Get 获取任务
func (r *Registry) Get(name string) Job {
	item, _ := r.items.Get(name)
	return item
}

// GetAll 获取所有任务
func (r *Registry) GetAll() []Job {
	return r.items.GetAll()
}

// GetNames 获取所有任务名称
func (r *Registry) GetNames() []string {
	return r.items.GetNames()
}

// Count 获取任务数量
func (r *Registry) Count() int {
	return r.items.Count()
}

// Exists 检查任务是否存在
func (r *Registry) Exists(name string) bool {
	return r.items.Exists(name)
}

// Register 全局注册任务
func Register(item Job) {
	if err := registry.Register(item); err != nil {
		flag.Errorf(pkg.Sprintf("Job [%s] 重复注册", item.Name()))
		os.Exit(1)
	}
}

// Get 全局获取任务
func Get(name string) Job {
	return registry.Get(name)
}

// GetAll 全局获取全部任务
func GetAll() []Job {
	return registry.GetAll()
}

// GetNames 全局获取全部任务名称
func GetNames() []string {
	return registry.GetNames()
}

// Count 全局统计任务数量
func Count() int {
	return registry.Count()
}

// Exists 全局检查任务是否存在
func Exists(name string) bool {
	return registry.Exists(name)
}
