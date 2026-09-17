package orm

import (
	"gin/config"
	"sync"

	"gorm.io/gorm"
)

// Manager 数据库连接管理器
type Manager struct {
	mu          sync.RWMutex
	conf        *config.Config
	defaultName string
	connections map[string]*gorm.DB
}

// NewManager 创建数据库连接管理器
func NewManager(conf *config.Config) *Manager {
	name := "mysql"
	if conf != nil && conf.Databases.Driver != "" {
		name = conf.Databases.Driver
	}

	return &Manager{
		conf:        conf,
		defaultName: name,
		connections: make(map[string]*gorm.DB),
	}
}

// Connection 获取指定数据库连接
func (m *Manager) Connection(name string) *gorm.DB {
	if name == "" {
		name = m.defaultName
	}

	m.mu.RLock()
	db := m.connections[name]
	m.mu.RUnlock()
	if db != nil {
		return db
	}

	db = Connection(name, m.conf)
	m.mu.Lock()
	m.connections[name] = db
	m.mu.Unlock()
	return db
}

// Reset 重置指定数据库连接
func (m *Manager) Reset(name string) *gorm.DB {
	if name == "" {
		name = m.defaultName
	}

	m.mu.Lock()
	delete(m.connections, name)
	m.mu.Unlock()

	db := ResetConnection(name)
	m.mu.Lock()
	m.connections[name] = db
	m.mu.Unlock()
	return db
}
