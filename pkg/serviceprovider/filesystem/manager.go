package filesystem

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"gin/config"

	"github.com/pkg6/gfs"
	"github.com/pkg6/gfs/bosfs"
	"github.com/pkg6/gfs/cosfs"
	"github.com/pkg6/gfs/kodofs"
	"github.com/pkg6/gfs/localfs"
	"github.com/pkg6/gfs/ossfs"
)

var (
	// ErrNotInitialized 文件系统服务未初始化
	ErrNotInitialized = errors.New("文件系统服务未初始化")
	// ErrDiskNotFound 文件磁盘不存在
	ErrDiskNotFound = errors.New("文件磁盘不存在")
	// ErrInvalidPath 文件路径不合法
	ErrInvalidPath = errors.New("文件路径不合法")
	// ErrUnsupportedSource 不支持的上传来源
	ErrUnsupportedSource = errors.New("不支持的上传来源")
)

// diskEntry 文件磁盘配置
type diskEntry struct {
	name        string
	adapter     gfs.IAdapter
	root        string
	displayRoot string
}

// Manager 文件系统管理器
type Manager struct {
	mu          sync.RWMutex
	defaultDisk string
	disks       map[string]*diskEntry
	order       []string
}

// NewManager 创建文件系统管理器
func NewManager(conf config.Filesystem) (*Manager, error) {
	manager := &Manager{
		defaultDisk: "local",
		disks:       make(map[string]*diskEntry),
	}

	manager.register(
		[]string{"local"},
		localfs.NewLocal(&localfs.Config{CDN: strings.TrimSpace(conf.Local.CDN)}),
		resolveLocalRoot(conf.Local.Root),
		resolveDisplayRoot(conf.Local.Root),
	)

	if ossConfigured(conf.OSS) {
		manager.register(
			[]string{"oss", "aliyun"},
			ossfs.NewOSS(&ossfs.Config{
				CDN:             strings.TrimSpace(conf.OSS.CDN),
				Bucket:          strings.TrimSpace(conf.OSS.Bucket),
				Endpoint:        strings.TrimSpace(conf.OSS.Endpoint),
				AccessKeyID:     strings.TrimSpace(conf.OSS.AccessKeyID),
				AccessKeySecret: strings.TrimSpace(conf.OSS.AccessKeySecret),
			}),
			"",
		)
	}

	if cosConfigured(conf.COS) {
		manager.register(
			[]string{"cos", "tencent"},
			cosfs.NewCOS(&cosfs.Config{
				CDN:       strings.TrimSpace(conf.COS.CDN),
				BucketURL: strings.TrimSpace(conf.COS.BucketURL),
				SecretID:  strings.TrimSpace(conf.COS.SecretID),
				SecretKey: strings.TrimSpace(conf.COS.SecretKey),
			}),
			"",
		)
	}

	if bosConfigured(conf.BOS) {
		manager.register(
			[]string{"bos", "baidu"},
			bosfs.NewBOS(&bosfs.Config{
				CDN:              strings.TrimSpace(conf.BOS.CDN),
				Ak:               strings.TrimSpace(conf.BOS.AK),
				Sk:               strings.TrimSpace(conf.BOS.SK),
				Endpoint:         strings.TrimSpace(conf.BOS.Endpoint),
				Bucket:           strings.TrimSpace(conf.BOS.Bucket),
				RedirectDisabled: conf.BOS.RedirectDisabled,
			}),
			"",
		)
	}

	if kodoConfigured(conf.KODO) {
		manager.register(
			[]string{"kodo", "qiniu"},
			kodofs.NewKoDo(&kodofs.Config{
				CDN:       strings.TrimSpace(conf.KODO.CDN),
				AccessKey: strings.TrimSpace(conf.KODO.AccessKey),
				SecretKey: strings.TrimSpace(conf.KODO.SecretKey),
				Bucket:    strings.TrimSpace(conf.KODO.Bucket),
			}),
			"",
		)
	}

	if name := normalizeDiskName(conf.Default); name != "" {
		manager.defaultDisk = name
	}
	if _, ok := manager.disks[manager.defaultDisk]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrDiskNotFound, manager.defaultDisk)
	}

	return manager, nil
}

// register 注册文件磁盘
func (m *Manager) register(names []string, adapter gfs.IAdapter, root string, displayRoot ...string) {
	if m == nil || adapter == nil || len(names) == 0 {
		return
	}

	name := normalizeDiskName(names[0])
	entry := &diskEntry{
		name:        name,
		adapter:     adapter,
		root:        root,
		displayRoot: firstDisplayRoot(displayRoot...),
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.order = append(m.order, name)
	for _, alias := range names {
		m.disks[normalizeDiskName(alias)] = entry
	}
}

// Disk 获取文件磁盘
func (m *Manager) Disk(name ...string) *Disk {
	if m == nil {
		return &Disk{err: ErrNotInitialized}
	}

	diskName := m.defaultDisk
	if len(name) > 0 && strings.TrimSpace(name[0]) != "" {
		diskName = normalizeDiskName(name[0])
	}

	m.mu.RLock()
	entry := m.disks[diskName]
	m.mu.RUnlock()
	if entry == nil {
		return &Disk{
			name: diskName,
			err:  fmt.Errorf("%w: %s", ErrDiskNotFound, diskName),
		}
	}

	return &Disk{
		name:        entry.name,
		adapter:     entry.adapter,
		root:        entry.root,
		displayRoot: entry.displayRoot,
	}
}

// Default 获取默认文件磁盘
func (m *Manager) Default() *Disk {
	return m.Disk()
}

// Has 判断文件磁盘是否存在
func (m *Manager) Has(name string) bool {
	if m == nil {
		return false
	}

	m.mu.RLock()
	_, ok := m.disks[normalizeDiskName(name)]
	m.mu.RUnlock()

	return ok
}

// Disks 获取已注册文件磁盘
func (m *Manager) Disks() []string {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]string(nil), m.order...)
}

// normalizeDiskName 标准化文件磁盘名称
func normalizeDiskName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "local":
		return "local"
	case "oss", "aliyun", "ali":
		return "oss"
	case "cos", "tencent":
		return "cos"
	case "bos", "baidu":
		return "bos"
	case "kodo", "qiniu":
		return "kodo"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

// resolveLocalRoot 解析本地存储根目录
func resolveLocalRoot(root string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		root = filepath.Join("storage", "uploads")
	}
	if !filepath.IsAbs(root) {
		root = filepath.Join(config.RootPath(), root)
	}
	return filepath.Clean(root)
}

// resolveDisplayRoot 解析返回对象使用的根目录
func resolveDisplayRoot(root string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		root = filepath.Join("storage", "uploads")
	}
	return path.Clean(strings.ReplaceAll(root, "\\", "/"))
}

// firstDisplayRoot 获取首个返回根目录
func firstDisplayRoot(values ...string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

// ossConfigured 阿里云OSS是否已配置
func ossConfigured(conf config.OSSFilesystem) bool {
	return conf.Bucket != "" && conf.AccessKeyID != "" && conf.AccessKeySecret != ""
}

// cosConfigured 腾讯云COS是否已配置
func cosConfigured(conf config.COSFilesystem) bool {
	return conf.BucketURL != "" && conf.SecretID != "" && conf.SecretKey != ""
}

// bosConfigured 百度云BOS是否已配置
func bosConfigured(conf config.BOSFilesystem) bool {
	return conf.Bucket != "" && conf.AK != "" && conf.SK != ""
}

// kodoConfigured 七牛云KODO是否已配置
func kodoConfigured(conf config.KODOFilesystem) bool {
	return conf.Bucket != "" && conf.AccessKey != "" && conf.SecretKey != ""
}
