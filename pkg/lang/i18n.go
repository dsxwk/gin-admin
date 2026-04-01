package lang

import (
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	Bundle     *i18n.Bundle
	Localizers = map[string]*i18n.Localizer{}
	once       sync.Once
	log        *logger.Logger
)

// LoadLang 初始化翻译
func LoadLang(conf *config.Config, logger *logger.Logger) {
	once.Do(func() {
		log = logger
		Bundle = i18n.NewBundle(language.Chinese)
		Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		baseDir := conf.I18n.Dir
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			log.Info(pkg.Sprintf("i18n baseDir not found: %s", baseDir))
			return
		}

		langs := strings.Split(conf.I18n.Lang, ",")

		// 遍历语言目录
		for _, langCode := range langs {
			langDir := filepath.Join(baseDir, langCode)
			loadLangDir(langCode, langDir)
		}

		// 初始化Localizer
		for _, langCode := range langs {
			Localizers[langCode] = i18n.NewLocalizer(Bundle, langCode)
		}

		log.Info(pkg.Sprintf("翻译服务加载成功,支持语言: %s", conf.I18n.Lang))
	})
}

// loadLangDir 递归加载指定语言目录下的所有翻译文件
func loadLangDir(langCode, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Info(pkg.Sprintf("遍历翻译目录失败: %v", err))
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			log.Info(pkg.Sprintf("读取翻译文件失败: %v", err))
			return nil
		}

		// 模拟路径格式如zh.json/en.yaml,让go-i18n能识别语言
		virtualFileName := fmt.Sprintf("%s%s", langCode, ext)
		_, err = Bundle.ParseMessageFileBytes(data, virtualFileName)
		if err != nil {
			log.Info(pkg.Sprintf("解析翻译文件失败: %v", err))
		}

		return nil
	})
	if err != nil {
		log.Info(pkg.Sprintf("加载翻译目录失败: %v", err))
	}
}

// T 翻译
func T(ctx context.Context, messageID string, data map[string]interface{}) string {
	langCode := getLangFromContext(ctx)

	localizer, ok := Localizers[langCode]
	if !ok {
		localizer = Localizers["zh"]
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		log.Debug(pkg.Sprintf("缺少翻译: %s (%s)", messageID, langCode))
		return messageID
	}
	return msg
}

// getLangFromContext 从上下文获取语言
func getLangFromContext(ctx context.Context) string {
	if ctx == nil {
		return "zh"
	}
	if v := ctx.Value(ctxkey.LangKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return "zh"
}

// GetLocalizer 获取指定语言
func GetLocalizer(langCode string) *i18n.Localizer {
	if localizer, ok := Localizers[langCode]; ok {
		return localizer
	}
	return Localizers["zh"]
}

// GetBundle 获取翻译包
func GetBundle() *i18n.Bundle {
	return Bundle
}

// IsLoaded 检查翻译是否已加载
func IsLoaded() bool {
	return Bundle != nil && len(Localizers) > 0
}
