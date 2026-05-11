package facade

import (
	"sync"
)

// Manager 门面管理-所有注册的服务实例
type Manager struct {
	mu        sync.RWMutex
	instances map[string]any // 存储服务实例
}

var (
	globalManager *Manager
	once          sync.Once
)

// init 自动创建
func init() {
	Init()
}

// Init 初始化
func Init() {
	once.Do(func() {
		globalManager = &Manager{
			instances: make(map[string]any),
		}
	})
}

// GetManager 获取门面管理器
func GetManager() *Manager {
	if globalManager == nil {
		Init()
	}
	return globalManager
}

// Register 注册服务实例到门面(泛型)
// 使用示例: facade.Register[*gorm.DB]("db", orm.Connection())
func Register[T any](name string, instance T) {
	mgr := GetManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.instances[name] = instance
}

// Get 从门面获取服务实例(泛型)
// 使用示例: db := facade.Get[*gorm.DB]("db")
func Get[T any](name string) T {
	mgr := GetManager()
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	if instance, ok := mgr.instances[name]; ok {
		if typed, ok := instance.(T); ok {
			return typed
		}
	}
	var zero T
	return zero
}

// Has 检查服务是否存在
func Has(name string) bool {
	mgr := GetManager()
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	_, ok := mgr.instances[name]
	return ok
}
