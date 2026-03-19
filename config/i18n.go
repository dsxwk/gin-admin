package config

// I18n 翻译
type I18n struct {
	Dir  string `mapstructure:"dir" yaml:"dir"`   // 翻译文件目录
	Lang string `mapstructure:"lang" yaml:"lang"` // 默认语言,多个语言用逗号分隔
}
