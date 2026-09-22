package lang

import (
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/config"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/goccy/go-json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var defaultService atomic.Pointer[Service]

// langState 翻译运行状态
type langState struct {
	bundle      *i18n.Bundle
	localizers  map[string]*i18n.Localizer
	defaultLang string
}

// Service 翻译服务
type Service struct {
	state atomic.Pointer[langState]
}

// Load 加载翻译文件
func (s *Service) Load(conf *config.Config) error {
	if s == nil {
		return fmt.Errorf("翻译服务未初始化")
	}
	if conf == nil {
		return fmt.Errorf("翻译配置未初始化")
	}
	if conf.I18n.Dir == "" {
		return fmt.Errorf("翻译目录未配置")
	}

	baseDir := conf.I18n.Dir
	if !filepath.IsAbs(baseDir) {
		baseDir = filepath.Join(config.GetRootPath(), baseDir)
	}
	if _, err := os.Stat(baseDir); err != nil {
		return fmt.Errorf("翻译目录不可用: %s: %w", baseDir, err)
	}

	languageCodes := parseLanguages(conf.I18n.Lang)
	if len(languageCodes) == 0 {
		return fmt.Errorf("翻译语言未配置")
	}

	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	localizers := make(map[string]*i18n.Localizer, len(languageCodes))
	for _, languageCode := range languageCodes {
		langDir := filepath.Join(baseDir, languageCode)
		if err := loadLangDir(bundle, languageCode, langDir); err != nil {
			return fmt.Errorf("加载%s翻译失败: %w", languageCode, err)
		}
		localizers[languageCode] = i18n.NewLocalizer(bundle, languageCode)
	}

	s.state.Store(&langState{
		bundle:      bundle,
		localizers:  localizers,
		defaultLang: languageCodes[0],
	})
	defaultService.Store(s)

	return nil
}

// parseLanguages 解析语言列表
func parseLanguages(value string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	for item := range strings.SplitSeq(value, ",") {
		languageCode := strings.TrimSpace(item)
		if languageCode == "" {
			continue
		}
		if _, ok := seen[languageCode]; ok {
			continue
		}

		seen[languageCode] = struct{}{}
		result = append(result, languageCode)
	}

	return result
}

// loadLangDir 递归加载指定语言目录下的所有翻译文件
func loadLangDir(bundle *i18n.Bundle, langCode, dir string) error {
	return filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取翻译文件失败: %s: %w", path, err)
		}

		// 模拟路径格式如zh.json/en.yaml,让go-i18n能识别语言
		virtualFileName := langCode + ext
		if _, err = bundle.ParseMessageFileBytes(data, virtualFileName); err != nil {
			return fmt.Errorf("解析翻译文件失败: %s: %w", path, err)
		}

		return nil
	})
}

// New 创建翻译服务
func New() *Service {
	return &Service{}
}

// Trans 翻译
func Trans(ctx context.Context, messageID string, data map[string]any) string {
	service := defaultService.Load()
	if service == nil {
		return messageID
	}
	return service.Trans(ctx, messageID, data)
}

// Trans 翻译
func (s *Service) Trans(ctx context.Context, messageID string, data map[string]any) string {
	state := s.currentState()
	if state == nil {
		return messageID
	}

	localizer := selectLocalizer(state, getLangFromContext(ctx))
	if localizer == nil {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}

	return msg
}

// Localizer 获取指定语言
func (s *Service) Localizer(langCode string) *i18n.Localizer {
	return selectLocalizer(s.currentState(), langCode)
}

// Bundle 获取翻译包
func (s *Service) Bundle() *i18n.Bundle {
	state := s.currentState()
	if state == nil {
		return nil
	}
	return state.bundle
}

// IsLoaded 检查翻译是否已加载
func (s *Service) IsLoaded() bool {
	state := s.currentState()
	return state != nil && state.bundle != nil && len(state.localizers) > 0
}

// currentState 获取当前翻译状态
func (s *Service) currentState() *langState {
	if s == nil {
		return nil
	}
	return s.state.Load()
}

// selectLocalizer 选择语言翻译器
func selectLocalizer(state *langState, langCode string) *i18n.Localizer {
	if state == nil {
		return nil
	}

	if langCode == "" {
		langCode = state.defaultLang
	}
	if localizer, ok := state.localizers[langCode]; ok {
		return localizer
	}

	return state.localizers[state.defaultLang]
}

// getLangFromContext 从上下文获取语言
func getLangFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value := ctx.Value(ctxkey.LangKey); value != nil {
		if languageCode, ok := value.(string); ok {
			return strings.TrimSpace(languageCode)
		}
	}
	return ""
}
