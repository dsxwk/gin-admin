package config

import "time"

// Cache 缓存
type Cache struct {
	Driver string `mapstructure:"driver" yaml:"driver"`
	Redis  Redis  `mapstructure:"redis" yaml:"redis"`
	Memory Memory `mapstructure:"memory" yaml:"memory"`
	Disk   Disk   `mapstructure:"disk" yaml:"disk"`
}

// Redis 数据库
type Redis struct {
	Address  string `mapstructure:"address" yaml:"address"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

// Memory 内存缓存
type Memory struct {
	DefaultExpire   time.Duration `mapstructure:"default-expire" yaml:"default-expire"`
	CleanupInterval time.Duration `mapstructure:"cleanup-interval" yaml:"cleanup-interval"`
}

// Disk 磁盘缓存
type Disk struct {
	Path string `mapstructure:"path" yaml:"path"`
}
