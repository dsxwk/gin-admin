package facade

import (
	"context"
	"fmt"
	"gin/app/errcode"
	"gin/app/model"
	"gin/pkg/serviceprovider/filesystem"
	"io"
	"path"
	"time"

	"github.com/pkg6/gfs"
	"gorm.io/gorm"
)

// File 文件门面
func File() *FileFacade {
	return &FileFacade{
		manager: App().File(),
		ctx:     context.Background(),
	}
}

// FileFacade 文件门面
type FileFacade struct {
	manager    *filesystem.Manager
	ctx        context.Context
	uploaderId int64
	groupId    int64
}

// WithContext 设置上下文
func (f *FileFacade) WithContext(ctx context.Context) *FileFacade {
	if f == nil {
		return f
	}
	result := *f
	if ctx != nil {
		result.ctx = ctx
	}
	return &result
}

// WithUploader 设置上传者ID
func (f *FileFacade) WithUploader(uploaderId int64) *FileFacade {
	if f == nil {
		return f
	}
	result := *f
	result.uploaderId = uploaderId
	return &result
}

// WithGroup 设置文件分组ID
func (f *FileFacade) WithGroup(groupId int64) *FileFacade {
	if f == nil {
		return f
	}
	result := *f
	result.groupId = groupId
	return &result
}

// Disk 获取文件磁盘
func (f *FileFacade) Disk(name ...string) *FileDisk {
	if f == nil {
		return &FileDisk{}
	}
	return &FileDisk{
		disk:       f.manager.Disk(name...),
		ctx:        f.ctx,
		uploaderId: f.uploaderId,
		groupId:    f.groupId,
	}
}

// Has 判断文件磁盘是否存在
func (f *FileFacade) Has(name string) bool {
	if f == nil || f.manager == nil {
		return false
	}
	return f.manager.Has(name)
}

// Disks 获取全部文件磁盘
func (f *FileFacade) Disks() []string {
	if f == nil || f.manager == nil {
		return nil
	}
	return f.manager.Disks()
}

// FileDisk 文件磁盘
type FileDisk struct {
	disk       *filesystem.Disk
	ctx        context.Context
	uploaderId int64
	groupId    int64
}

// Name 获取文件磁盘名称
func (d *FileDisk) Name() string {
	if d == nil || d.disk == nil {
		return ""
	}
	return d.disk.Name()
}

// Adapter 获取底层文件适配器
func (d *FileDisk) Adapter() gfs.IAdapter {
	if d == nil || d.disk == nil {
		return nil
	}
	return d.disk.Adapter()
}

// Upload 上传文件
func (d *FileDisk) Upload(source any, objectPath string) (*filesystem.UploadResult, error) {
	if d == nil || d.disk == nil {
		return nil, (errcode.FilesystemErrCode{}).NotInitialized()
	}
	result, err := d.disk.Upload(source, objectPath)
	if err != nil {
		return nil, err
	}
	if err = d.record(result); err != nil {
		_ = d.disk.Delete(result.ObjectKey)
		return nil, err
	}
	return result, nil
}

// UploadReader 上传文件流
func (d *FileDisk) UploadReader(reader io.Reader, objectPath string) (*filesystem.UploadResult, error) {
	return d.UploadReaderWithSize(reader, objectPath, 0)
}

// UploadReaderWithSize 上传文件流并指定大小
func (d *FileDisk) UploadReaderWithSize(reader io.Reader, objectPath string, size int64) (*filesystem.UploadResult, error) {
	if d == nil || d.disk == nil {
		return nil, (errcode.FilesystemErrCode{}).NotInitialized()
	}
	result, err := d.disk.UploadReaderWithSize(reader, objectPath, size)
	if err != nil {
		return nil, err
	}
	if err = d.record(result); err != nil {
		_ = d.disk.Delete(result.ObjectKey)
		return nil, err
	}
	return result, nil
}

// Read 读取文件
func (d *FileDisk) Read(objectPath string) ([]byte, error) {
	if d == nil || d.disk == nil {
		return nil, (errcode.FilesystemErrCode{}).NotInitialized()
	}
	return d.disk.Read(objectPath)
}

// Delete 删除文件和文件记录
func (d *FileDisk) Delete(objectPath string) error {
	if d == nil || d.disk == nil {
		return (errcode.FilesystemErrCode{}).NotInitialized()
	}
	objectKey, err := d.disk.ObjectKey(objectPath)
	if err != nil {
		return err
	}

	db := DB("mysql")
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.WithContext(d.context()).Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("storage = ? AND object_key = ?", d.disk.Name(), objectKey).
			Delete(&model.File{}).Error; err != nil {
			return fmt.Errorf("删除文件记录失败:%w", err)
		}

		exists, err := d.disk.Exists(objectKey)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		return d.disk.Delete(objectKey)
	})
}

// Exists 判断文件是否存在
func (d *FileDisk) Exists(objectPath string) (bool, error) {
	if d == nil || d.disk == nil {
		return false, (errcode.FilesystemErrCode{}).NotInitialized()
	}
	return d.disk.Exists(objectPath)
}

// URL 获取文件访问地址
func (d *FileDisk) URL(objectPath string) (string, error) {
	if d == nil || d.disk == nil {
		return "", (errcode.FilesystemErrCode{}).NotInitialized()
	}
	return d.disk.URL(objectPath)
}

// record 保存上传文件记录
func (d *FileDisk) record(result *filesystem.UploadResult) error {
	if result == nil {
		return fmt.Errorf("文件上传结果为空")
	}
	db := DB("mysql")
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	objectKey := result.ObjectKey
	if objectKey == "" {
		objectKey = path.Base(result.Path)
	}
	originName := result.OriginName
	if originName == "" {
		originName = path.Base(objectKey)
	}

	now := model.DateTime(time.Now())
	record := model.File{
		GroupId:    d.groupId,
		Storage:    result.Disk,
		OriginName: originName,
		ObjectKey:  objectKey,
		Url:        result.URL,
		MimeType:   result.MimeType,
		Size:       result.Size,
		UploaderId: d.uploaderId,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if err := db.WithContext(d.context()).Create(&record).Error; err != nil {
		return fmt.Errorf("保存文件记录失败:%w", err)
	}
	return nil
}

// context 获取文件操作上下文
func (d *FileDisk) context() context.Context {
	if d == nil || d.ctx == nil {
		return context.Background()
	}
	return d.ctx
}
