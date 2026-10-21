package orm

import (
	"fmt"
	"gin/config"
	servicelogger "gin/pkg/serviceprovider/logger"
	"gin/pkg/terminal"
	"sync"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Manager 数据库连接管理器
type Manager struct {
	mu          sync.RWMutex
	conf        *config.Config
	log         *servicelogger.Logger
	gormLogger  gormlogger.Interface
	defaultName string
	connections map[string]*gorm.DB
}

// NewManager 创建数据库连接管理器
func NewManager(conf *config.Config, logs ...*servicelogger.Logger) *Manager {
	name := "mysql"
	if conf != nil && conf.Databases.Default != "" {
		name = conf.Databases.Default
	}

	var log *servicelogger.Logger
	if len(logs) > 0 {
		log = logs[0]
	}

	return &Manager{
		conf:        conf,
		log:         log,
		gormLogger:  newGormLogger(conf),
		defaultName: name,
		connections: make(map[string]*gorm.DB),
	}
}

// Connect 创建并获取指定数据库连接
func (m *Manager) Connect(name string) (*gorm.DB, error) {
	if m == nil {
		return nil, fmt.Errorf("orm manager not initialized")
	}
	if name == "" {
		name = m.defaultName
	}

	m.mu.RLock()
	db := m.connections[name]
	m.mu.RUnlock()
	if db != nil {
		return db, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if db = m.connections[name]; db != nil {
		return db, nil
	}

	db, err := openConnection(name, m.conf, m.gormLogger)
	if err != nil {
		return nil, fmt.Errorf("%s数据库连接失败: %w", name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%s数据库连接池配置失败: %w", name, err)
	}

	maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime := connectionPoolSettings(m.conf)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	if err = sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("%s数据库连接ping失败: %w", name, err)
	}

	sqlCallback(db, m.conf, m.log)
	m.connections[name] = db
	return db, nil
}

// Connection 获取指定数据库连接,兼容原调用方式
func (m *Manager) Connection(name string) *gorm.DB {
	db, err := m.Connect(name)
	if err != nil {
		terminal.Errorf("%v", err)
		return nil
	}
	return db
}

// Reset 重置指定数据库连接
func (m *Manager) Reset(name string) *gorm.DB {
	db, err := m.Reconnect(name)
	if err != nil {
		terminal.Errorf("%v", err)
		return nil
	}
	return db
}

// Reconnect 关闭并重建指定数据库连接
func (m *Manager) Reconnect(name string) (*gorm.DB, error) {
	if m == nil {
		return nil, fmt.Errorf("orm manager not initialized")
	}
	if name == "" {
		name = m.defaultName
	}

	m.mu.Lock()
	if db := m.connections[name]; db != nil {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		delete(m.connections, name)
	}
	m.mu.Unlock()

	return m.Connect(name)
}
