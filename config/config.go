package config

import (
	"fmt"
	"gin/common/flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config 配置
type Config struct {
	App       App       `mapstructure:"app" yaml:"app"`
	Databases Databases `mapstructure:"databases" yaml:"databases"`
	Cors      Cors      `mapstructure:"cors" yaml:"cors"`
	Jwt       Jwt       `mapstructure:"jwt" yaml:"jwt"`
	Log       Log       `mapstructure:"log" yaml:"log"`
	Cache     Cache     `mapstructure:"cache" yaml:"cache"`
	I18n      I18n      `mapstructure:"i18n" yaml:"i18n"`
	Kafka     Kafka     `mapstructure:"kafka" yaml:"kafka"`
	Rabbitmq  Rabbitmq  `mapstructure:"rabbitmq" yaml:"rabbitmq"`
}

var (
	Conf     *Config
	vp       *viper.Viper
	confOnce sync.Once
)

func NewConfig() *Config {
	confOnce.Do(func() {
		v := viper.New()

		// 默认配置文件目录为根目录
		configDir := GetRootPath()
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

		// 加载对应环境的配置文件，如 dev.config.yaml
		configFile := filepath.Join(configDir, fmt.Sprintf("%s.config.yaml", env))
		if _, err := os.Stat(configFile); err == nil {
			v.SetConfigFile(configFile)
			if err = v.MergeInConfig(); err != nil {
				flag.Errorf("合并环境配置失败: %v", err)
				os.Exit(1)
			}
			flag.Successf("已加载环境配置文件: %s", configFile)
		} else {
			flag.Warningf("未找到环境配置文件: %s，使用默认配置\n", configFile)
		}

		// 自动映射到结构体
		cfg := &Config{}
		if err := v.Unmarshal(cfg); err != nil {
			flag.Errorf("解析配置文件失败: %v", err)
			os.Exit(1)
		}

		// 监听配置变化
		v.WatchConfig()

		var lastEventTime int64
		v.OnConfigChange(func(e fsnotify.Event) {
			if e.Op&fsnotify.Write != fsnotify.Write {
				return
			}

			now := time.Now().UnixNano()
			// 如果两次事件间隔小于200ms则忽略
			if now-lastEventTime < 200*1e6 {
				return
			}
			lastEventTime = now

			flag.Infof("配置文件修改: %s\n", e.Name)
			if err := v.Unmarshal(cfg); err != nil {
				flag.Errorf("配置热更新失败: %v", err)
				os.Exit(1)
			}
		})

		Conf = cfg
		vp = v
	})

	return Conf
}

// GetRootPath 获取项目根路径
func GetRootPath() string {
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
func Get(key string) interface{} {
	return vp.Get(key)
}

// GetString 获取字符串
func GetString(key string) string {
	return vp.GetString(key)
}

// GetInt 获取整数
func GetInt(key string) int {
	return vp.GetInt(key)
}

// GetBool 获取布尔值
func GetBool(key string) bool {
	return vp.GetBool(key)
}
