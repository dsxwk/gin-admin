package config

// Filesystem 文件系统配置
type Filesystem struct {
	Default string          `mapstructure:"default" yaml:"default"` // 默认磁盘
	Local   LocalFilesystem `mapstructure:"local" yaml:"local"`     // 本地磁盘
	OSS     OSSFilesystem   `mapstructure:"oss" yaml:"oss"`         // 阿里云OSS
	COS     COSFilesystem   `mapstructure:"cos" yaml:"cos"`         // 腾讯云COS
	BOS     BOSFilesystem   `mapstructure:"bos" yaml:"bos"`         // 百度云BOS
	KODO    KODOFilesystem  `mapstructure:"kodo" yaml:"kodo"`       // 七牛云KODO
}

// LocalFilesystem 本地文件系统配置
type LocalFilesystem struct {
	Root string `mapstructure:"root" yaml:"root"` // 本地存储根目录
	CDN  string `mapstructure:"cdn" yaml:"cdn"`   // 访问域名
}

// OSSFilesystem 阿里云OSS配置
type OSSFilesystem struct {
	CDN             string `mapstructure:"cdn" yaml:"cdn"`                             // 访问域名
	Bucket          string `mapstructure:"bucket" yaml:"bucket"`                       // 存储桶
	Endpoint        string `mapstructure:"endpoint" yaml:"endpoint"`                   // 地域节点
	AccessKeyID     string `mapstructure:"access-key-id" yaml:"access-key-id"`         // AccessKeyID
	AccessKeySecret string `mapstructure:"access-key-secret" yaml:"access-key-secret"` // AccessKeySecret
}

// COSFilesystem 腾讯云COS配置
type COSFilesystem struct {
	CDN       string `mapstructure:"cdn" yaml:"cdn"`               // 访问域名
	BucketURL string `mapstructure:"bucket-url" yaml:"bucket-url"` // 存储桶地址
	SecretID  string `mapstructure:"secret-id" yaml:"secret-id"`   // SecretID
	SecretKey string `mapstructure:"secret-key" yaml:"secret-key"` // SecretKey
}

// BOSFilesystem 百度云BOS配置
type BOSFilesystem struct {
	CDN              string `mapstructure:"cdn" yaml:"cdn"`                             // 访问域名
	AK               string `mapstructure:"ak" yaml:"ak"`                               // AccessKey
	SK               string `mapstructure:"sk" yaml:"sk"`                               // SecretKey
	Endpoint         string `mapstructure:"endpoint" yaml:"endpoint"`                   // 地域节点
	Bucket           string `mapstructure:"bucket" yaml:"bucket"`                       // 存储桶
	RedirectDisabled bool   `mapstructure:"redirect-disabled" yaml:"redirect-disabled"` // 是否禁用重定向
}

// KODOFilesystem 七牛云KODO配置
type KODOFilesystem struct {
	CDN       string `mapstructure:"cdn" yaml:"cdn"`               // 访问域名
	AccessKey string `mapstructure:"access-key" yaml:"access-key"` // AccessKey
	SecretKey string `mapstructure:"secret-key" yaml:"secret-key"` // SecretKey
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`         // 存储桶
}
