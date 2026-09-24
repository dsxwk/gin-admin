package cache

import (
	"context"
	"errors"
	"fmt"
	"gin/config"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/serviceprovider/logger"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// DiskCache 磁盘缓存
type DiskCache struct {
	db  *badger.DB
	ctx context.Context
}

var diskCache *CacheProxy

func NewDiskCache(conf *config.Config) *CacheProxy {
	if diskCache != nil {
		return diskCache
	}
	path := conf.Cache.Disk.Path
	if path != "" && !filepath.IsAbs(path) {
		if rootPath := config.RootPath(); rootPath != "" {
			path = filepath.Join(rootPath, path)
		}
	}
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		logger.NewLogger(conf).Error(fmt.Sprintf("init disk cache failed: %s", err.Error()))
	}
	disk := &DiskCache{db: db}

	diskCache = NewCacheProxy("disk", disk, eventbus.Default(), nil)
	return diskCache
}

func (d *DiskCache) WithContext(ctx context.Context) *DiskCache {
	return &DiskCache{
		db:  d.db,
		ctx: ctx,
	}
}

func (d *DiskCache) Set(key string, value any, expire time.Duration) error {
	data, err := encodeCacheValue(value)
	if err != nil {
		return err
	}

	return d.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), data)
		if expire > 0 {
			e = e.WithTTL(expire)
		}
		return txn.SetEntry(e)
	})
}

func (d *DiskCache) Get(key string) (any, bool) {
	data, _, ok := d.load(key)
	if !ok {
		return nil, false
	}

	// 缓存只存储JSON,解析失败按未命中处理
	value, err := decodeCacheValue(data)
	if err != nil {
		return nil, false
	}

	return value, true
}

func (d *DiskCache) Delete(key string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (d *DiskCache) Expire(key string) (any, time.Time, bool, error) {
	data, expireTime, ok := d.load(key)
	if !ok {
		return nil, time.Time{}, false, errors.New("cache key not found")
	}

	value, err := decodeCacheValue(data)
	if err != nil {
		return nil, time.Time{}, false, err
	}

	return value, expireTime, true, nil
}

// load 读取缓存原始数据和过期时间
func (d *DiskCache) load(key string) ([]byte, time.Time, bool) {
	var (
		data   []byte
		expire time.Time
	)

	// Badger 不直接支持获取剩余ttl,只能判断是否存在
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		if ttl := item.ExpiresAt(); ttl > 0 {
			expire = time.Unix(int64(ttl), 0)
		}
		return nil
	})
	if err != nil {
		return nil, time.Time{}, false
	}

	return data, expire, true
}

// Close 关闭磁盘缓存
func (d *DiskCache) Close() error {
	if d == nil || d.db == nil {
		return nil
	}

	return d.db.Close()
}
