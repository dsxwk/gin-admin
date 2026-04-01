package facade

import (
	"sync"
)

// Manager 门面管理-所有注册的服务实例
type Manager struct {
	mu        sync.RWMutex
	instances map[string]interface{} // 存储服务实例
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
			instances: make(map[string]interface{}),
		}
	})
}

// Register 注册服务实例到门面
// 使用示例: facade.Register("db", orm.Connection())
func Register(name string, instance interface{}) {
	mgr := GetManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.instances[name] = instance
}

// Get 从门面获取服务实例
// 使用示例: db := facade.Get("db").(*gorm.DB)
func Get(name string) interface{} {
	mgr := GetManager()
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	if instance, ok := mgr.instances[name]; ok {
		return instance
	}
	return nil
}

// Has 检查服务是否存在
func Has(name string) bool {
	mgr := GetManager()
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	_, ok := mgr.instances[name]
	return ok
}

// GetManager 获取门面管理
func GetManager() *Manager {
	if globalManager == nil {
		Init()
	}
	return globalManager
}
