package job

import (
	"gin/common/flag"
	"gin/pkg"
	"os"
	"sync"
)

// Job 任务接口
type Job interface {
	// Name 任务名称(唯一标识)
	Name() string
	// Description 任务描述
	Description() string
	// Connection 任务连接, 默认 "redis",可选 "sync","kafka","rabbitmq"
	Connection() string
	// Retry 重试次数, 默认3
	Retry() int
	// Delay 重试间隔时间(毫秒), 默认1000
	Delay() int64
	// NewPayload 返回空payload指针,消费层用于json.Unmarshal
	NewPayload() any
	// Handle 业务处理逻辑,payload已是反序列化后的结构体
	Handle(payload any) error
}

// Registry 任务注册表
type Registry struct {
	items map[string]Job
	mu    sync.RWMutex
}

var registry = &Registry{
	items: make(map[string]Job),
}

// Register 注册任务
func Register(job Job) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	name := job.Name()
	if _, exists := registry.items[name]; exists {
		flag.Errorf(pkg.Sprintf("Job [%s] 重复注册", name))
		os.Exit(1)
	}
	registry.items[name] = job
}

// Get 获取任务
func Get(name string) Job {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return registry.items[name]
}

// GetAll 获取所有任务
func GetAll() []Job {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	jobs := make([]Job, 0, len(registry.items))
	for _, j := range registry.items {
		jobs = append(jobs, j)
	}
	return jobs
}

// GetNames 获取所有任务名称
func GetNames() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	names := make([]string, 0, len(registry.items))
	for name := range registry.items {
		names = append(names, name)
	}
	return names
}

// Count 获取任务数量
func Count() int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return len(registry.items)
}

// Exists 检查任务是否存在
func Exists(name string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	_, exists := registry.items[name]
	return exists
}
