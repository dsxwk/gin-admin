package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg6/gfs"
	"github.com/pkg6/gfs/localfs"
)

// Disk 文件磁盘
type Disk struct {
	name        string
	adapter     gfs.IAdapter
	root        string
	displayRoot string
	err         error
}

// UploadResult 文件上传结果
type UploadResult struct {
	Disk     string `json:"disk"`     // 文件磁盘
	Path     string `json:"path"`     // 文件路径
	URL      string `json:"url"`      // 文件地址
	Size     int64  `json:"size"`     // 文件大小
	MimeType string `json:"mimeType"` // 文件类型
}

// Name 获取文件磁盘名称
func (d *Disk) Name() string {
	if d == nil {
		return ""
	}
	return d.name
}

// Adapter 获取底层文件适配器
func (d *Disk) Adapter() gfs.IAdapter {
	if d == nil {
		return nil
	}
	return d.adapter
}

// Upload 上传文件
func (d *Disk) Upload(source any, objectPath string) (*UploadResult, error) {
	reader, closer, size, err := openUploadSource(source)
	if err != nil {
		return nil, err
	}
	if closer != nil {
		defer func() {
			_ = closer.Close()
		}()
	}

	return d.UploadReaderWithSize(reader, objectPath, size)
}

// UploadReader 上传文件流
func (d *Disk) UploadReader(reader io.Reader, objectPath string) (*UploadResult, error) {
	return d.UploadReaderWithSize(reader, objectPath, 0)
}

// UploadReaderWithSize 上传文件流并指定文件大小
func (d *Disk) UploadReaderWithSize(reader io.Reader, objectPath string, size int64) (*UploadResult, error) {
	if err := d.ready(); err != nil {
		return nil, err
	}
	if reader == nil {
		return nil, ErrUnsupportedSource
	}

	objectPath, err := normalizeObjectPath(objectPath)
	if err != nil {
		return nil, err
	}
	objectPath = d.objectKey(objectPath)

	storedPath, err := d.storedPath(objectPath)
	if err != nil {
		return nil, err
	}

	if err = d.adapter.WriteReader(storedPath, reader); err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	result := &UploadResult{
		Disk: d.name,
		Path: d.resultPath(objectPath),
		Size: size,
	}

	if result.Size <= 0 {
		if storedSize, sizeErr := d.adapter.Size(storedPath); sizeErr == nil {
			result.Size = storedSize
		}
	}

	if mimeType, mimeErr := d.adapter.MimeType(storedPath); mimeErr == nil {
		result.MimeType = mimeType
	} else {
		result.MimeType = mime.TypeByExtension(path.Ext(objectPath))
	}

	if uri, urlErr := d.adapter.URL(objectPath); urlErr == nil && uri != nil {
		d.applyDisplayRoot(uri)
		result.URL = uri.String()
	}

	return result, nil
}

// Read 读取文件
func (d *Disk) Read(objectPath string) ([]byte, error) {
	if err := d.ready(); err != nil {
		return nil, err
	}

	objectPath, err := normalizeObjectPath(objectPath)
	if err != nil {
		return nil, err
	}
	objectPath = d.objectKey(objectPath)

	storedPath, err := d.storedPath(objectPath)
	if err != nil {
		return nil, err
	}

	return d.adapter.Read(storedPath)
}

// Delete 删除文件
func (d *Disk) Delete(objectPath string) error {
	if err := d.ready(); err != nil {
		return err
	}

	objectPath, err := normalizeObjectPath(objectPath)
	if err != nil {
		return err
	}
	objectPath = d.objectKey(objectPath)

	storedPath, err := d.storedPath(objectPath)
	if err != nil {
		return err
	}

	if _, err = d.adapter.Delete(storedPath); err != nil {
		return err
	}

	if _, ok := d.adapter.(*localfs.Adapter); ok {
		if _, statErr := os.Stat(storedPath); statErr == nil {
			return fmt.Errorf("删除文件失败: %s", objectPath)
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return statErr
		}
	}

	return nil
}

// Exists 判断文件是否存在
func (d *Disk) Exists(objectPath string) (bool, error) {
	if err := d.ready(); err != nil {
		return false, err
	}

	objectPath, err := normalizeObjectPath(objectPath)
	if err != nil {
		return false, err
	}
	objectPath = d.objectKey(objectPath)

	storedPath, err := d.storedPath(objectPath)
	if err != nil {
		return false, err
	}

	exists, err := d.adapter.Exist(storedPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return exists, err
}

// URL 获取文件访问地址
func (d *Disk) URL(objectPath string) (string, error) {
	if err := d.ready(); err != nil {
		return "", err
	}

	objectPath, err := normalizeObjectPath(objectPath)
	if err != nil {
		return "", err
	}
	objectPath = d.objectKey(objectPath)

	uri, err := d.adapter.URL(objectPath)
	if err != nil {
		return "", err
	}
	if uri == nil {
		return "", nil
	}
	d.applyDisplayRoot(uri)

	return uri.String(), nil
}

// ready 检查文件磁盘是否可用
func (d *Disk) ready() error {
	if d == nil {
		return ErrNotInitialized
	}
	if d.err != nil {
		return d.err
	}
	if d.adapter == nil {
		return ErrNotInitialized
	}
	return nil
}

// storedPath 获取底层存储路径
func (d *Disk) storedPath(objectPath string) (string, error) {
	if err := d.ready(); err != nil {
		return "", err
	}
	if d.root == "" {
		return objectPath, nil
	}
	return filepath.Join(d.root, filepath.FromSlash(objectPath)), nil
}

// objectKey 去除返回路径中的根目录
func (d *Disk) objectKey(objectPath string) string {
	if d == nil || d.displayRoot == "" {
		return objectPath
	}

	prefix := strings.TrimSuffix(d.displayRoot, "/") + "/"
	if after, ok := strings.CutPrefix(objectPath, prefix); ok {
		return after
	}

	return objectPath
}

// resultPath 获取返回对象路径
func (d *Disk) resultPath(objectPath string) string {
	if d == nil || d.displayRoot == "" {
		return objectPath
	}
	return path.Join(d.displayRoot, objectPath)
}

// applyDisplayRoot 为本地访问地址补齐根目录
func (d *Disk) applyDisplayRoot(uri *url.URL) {
	if d == nil || uri == nil || d.displayRoot == "" {
		return
	}

	rootPath := "/" + strings.Trim(d.displayRoot, "/")
	currentPath := "/" + strings.Trim(uri.Path, "/")
	if currentPath == rootPath || strings.HasPrefix(currentPath, rootPath+"/") {
		return
	}

	uri.Path = path.Join(rootPath, currentPath)
	uri.RawPath = ""
}

// openUploadSource 打开上传来源
func openUploadSource(source any) (io.Reader, io.Closer, int64, error) {
	switch value := source.(type) {
	case nil:
		return nil, nil, 0, ErrUnsupportedSource
	case []byte:
		return bytes.NewReader(value), nil, int64(len(value)), nil
	case *multipart.FileHeader:
		file, err := value.Open()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("打开上传文件失败: %w", err)
		}
		return file, file, value.Size, nil
	case *os.File:
		info, err := value.Stat()
		if err != nil {
			return value, nil, 0, nil
		}
		return value, nil, info.Size(), nil
	case string:
		file, err := os.Open(strings.TrimSpace(value))
		if err != nil {
			return nil, nil, 0, fmt.Errorf("打开上传文件失败: %w", err)
		}
		info, err := file.Stat()
		if err != nil {
			return file, file, 0, nil
		}
		return file, file, info.Size(), nil
	case io.Reader:
		return value, nil, 0, nil
	default:
		return nil, nil, 0, ErrUnsupportedSource
	}
}

// normalizeObjectPath 标准化文件路径
func normalizeObjectPath(value string) (string, error) {
	value = strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	value = strings.TrimPrefix(value, "/")
	value = path.Clean(value)

	if value == "" || value == "." || value == ".." || strings.HasPrefix(value, "../") {
		return "", ErrInvalidPath
	}

	return value, nil
}
