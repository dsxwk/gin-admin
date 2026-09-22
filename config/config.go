package config

import (
	"fmt"
	"gin/common/flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	App            App            `mapstructure:"app" yaml:"app"`
	Databases      Databases      `mapstructure:"databases" yaml:"databases"`
	Cors           Cors           `mapstructure:"cors" yaml:"cors"`
	Jwt            Jwt            `mapstructure:"jwt" yaml:"jwt"`
	Log            Log            `mapstructure:"log" yaml:"log"`
	Cache          Cache          `mapstructure:"cache" yaml:"cache"`
	I18n           I18n           `mapstructure:"i18n" yaml:"i18n"`
	Queue          Queue          `mapstructure:"queue" yaml:"queue"`
	OperatorRecord OperatorRecord `mapstructure:"operator-record" yaml:"operator-record"`
	Mcp            Mcp            `mapstructure:"mcp" yaml:"mcp"`
	Agent          Agent          `mapstructure:"agent" yaml:"agent"`
	Grpc           Grpc           `mapstructure:"grpc" yaml:"grpc"`
	Es             Es             `mapstructure:"es" yaml:"es"`
}

var (
	conf            atomic.Pointer[Config]
	onConfigUpdated atomic.Pointer[func(*Config)]
	vp              *viper.Viper
	confOnce        sync.Once
	mu              sync.RWMutex // 添加读写锁保证并发安全
)

// SetOnConfigUpdated 设置配置更新回调
func SetOnConfigUpdated(callback func(*Config)) {
	if callback == nil {
		onConfigUpdated.Store(nil)
		return
	}
	onConfigUpdated.Store(&callback)
}

func NewConfig() *Config {
	confOnce.Do(func() {
		v := viper.New()

		// 默认配置文件目录为根目录
		configDir := RootPath()
		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		// 允许使用环境变量覆盖
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// 读取主配置文件 config.yaml
		if err := v.ReadInConfig(); err != nil {
			flag.Errorf("读取配置文件失败: %v", err)
		}

		// 获取环境类型
		env := v.GetString("app.env")
		if env == "" {
			env = "dev"
		}

		// 加载对应环境的配置文件,如 dev.config.yaml
		configFile := filepath.Join(configDir, fmt.Sprintf("%s.config.yaml", env))
		if _, err := os.Stat(configFile); err == nil {
			v.SetConfigFile(configFile)
			if err = v.MergeInConfig(); err != nil {
				flag.Errorf("合并环境配置失败: %v", err)
				os.Exit(1)
			}
			flag.Infof("已加载环境配置文件: %s", configFile)
		} else {
			defaultConfigFile := filepath.Join(configDir, "config.yaml")
			flag.Warningf("未找到环境配置文件 %s,使用默认配置 %s\n", configFile, defaultConfigFile)
		}

		// 自动映射到结构体
		cfg, err := unmarshalConfig(v)
		if err != nil {
			flag.Errorf("解析配置文件失败: %v", err)
			os.Exit(1)
		}

		mu.Lock()
		vp = v
		conf.Store(cfg)
		mu.Unlock()

		v.WatchConfig()

		v.OnConfigChange(func(e fsnotify.Event) {
			if e.Op&fsnotify.Write != fsnotify.Write {
				return
			}

			flag.Infof("配置文件修改: %s", e.Name)

			mu.Lock()
			next, err := unmarshalConfig(v)
			mu.Unlock()
			if err != nil {
				flag.Errorf("配置热更新失败: %v", err)
				return
			}

			conf.Store(next)
			notifyConfigUpdated(next)
		})
	})

	return conf.Load()
}

// unmarshalConfig 解析配置副本
func unmarshalConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// notifyConfigUpdated 通知配置更新
func notifyConfigUpdated(cfg *Config) {
	callback := onConfigUpdated.Load()
	if callback != nil {
		(*callback)(cfg)
	}
}

// RootPath 获取项目根路径
func RootPath() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// Get 获取配置项
func (c *Config) Get(key string) any {
	mu.RLock()
	defer mu.RUnlock()

	if vp == nil {
		return nil
	}
	return vp.Get(key)
}

// GetString 获取字符串
func (c *Config) GetString(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if vp == nil {
		return ""
	}
	return vp.GetString(key)
}

// GetInt 获取整数
func (c *Config) GetInt(key string) int {
	mu.RLock()
	defer mu.RUnlock()

	if vp == nil {
		return 0
	}
	return vp.GetInt(key)
}

// GetBool 获取布尔值
func (c *Config) GetBool(key string) bool {
	mu.RLock()
	defer mu.RUnlock()

	if vp == nil {
		return false
	}
	return vp.GetBool(key)
}

// Set 更新配置项(通用方法)
func (c *Config) Set(key string, value any) error {
	mu.Lock()

	if vp == nil {
		mu.Unlock()
		return fmt.Errorf("viper 未初始化")
	}

	// 获取当前配置文件路径
	configFile := vp.ConfigFileUsed()
	if configFile == "" {
		mu.Unlock()
		return fmt.Errorf("未找到配置文件")
	}

	// 设置新值
	vp.Set(key, value)

	// 写回文件
	if err := vp.WriteConfig(); err != nil {
		mu.Unlock()
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	// 重新加载配置到结构体
	next, err := unmarshalConfig(vp)
	if err != nil {
		mu.Unlock()
		return fmt.Errorf("重新加载配置失败: %w", err)
	}
	conf.Store(next)
	mu.Unlock()

	flag.Infof("配置已更新: %s = %v", key, value)
	notifyConfigUpdated(next)
	return nil
}

// SetString 更新字符串配置
func (c *Config) SetString(key, value string) error {
	return c.Set(key, value)
}

// SetInt 更新整数配置
func (c *Config) SetInt(key string, value int) error {
	return c.Set(key, value)
}

// SetBool 更新布尔配置
func (c *Config) SetBool(key string, value bool) error {
	return c.Set(key, value)
}
