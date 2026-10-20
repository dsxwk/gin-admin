package service

import (
	"context"
	"errors"
	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"mime/multipart"
	"path"
	"strings"
	"time"
	"uuid"

	"gorm.io/gorm"
)

// FileService 文件服务
type FileService struct {
	base.Service
}

// List 文件列表
func (s *FileService) List(ctx context.Context, req request.File) (pageData request.PageData, err error) {
	var (
		m      model.File
		models []model.File
		db     = s.DB(ctx, &m)
	)

	db = s.Search(db, m, req.Search).Model(&m).Preload("Group")
	if req.GroupId != nil {
		db = db.Where("group_id = ?", *req.GroupId)
	}

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("id DESC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	} else {
		pageData.Page = req.Page
		pageData.PageSize = req.PageSize
		offset, limit := request.Pagination(req.Page, req.PageSize)

		err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	}

	return pageData, nil
}

// Upload 上传文件
func (s *FileService) Upload(
	ctx context.Context,
	uploaderId int64,
	groupId int64,
	storage string,
	fileHeader *multipart.FileHeader,
) (m model.File, err error) {
	if fileHeader == nil {
		return m, errcode.ArgsError().WithMsg("请上传file文件")
	}
	if groupId < 0 {
		return m, errcode.ArgsError().WithMsg("分组ID不合法")
	}
	if groupId > 0 {
		var count int64
		err = s.DB(ctx, &model.FileGroup{}).
			Model(&model.FileGroup{}).
			Where("id = ?", groupId).
			Count(&count).Error
		if err != nil {
			return m, err
		}
		if count == 0 {
			return m, errcode.ArgsError().WithMsg("文件分组不存在")
		}
	}

	result, err := facade.File().
		WithContext(ctx).
		WithUploader(uploaderId).
		WithGroup(groupId).
		Disk(storage).
		Upload(fileHeader, fileObjectPath(fileHeader.Filename))
	if err != nil {
		return m, err
	}

	err = s.DB(ctx, &m).
		Preload("Group").
		Where("storage = ? AND object_key = ?", result.Disk, result.ObjectKey).
		Order("id DESC").
		First(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}

// Delete 删除文件
func (s *FileService) Delete(ctx context.Context, id int64) error {
	var m model.File
	err := s.DB(ctx, &m).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errcode.NotFound().WithMsg("文件不存在")
	}
	if err != nil {
		return err
	}

	return facade.File().
		WithContext(ctx).
		Disk(m.Storage).
		Delete(m.ObjectKey)
}

// UpdateGroup 修改文件分组
func (s *FileService) UpdateGroup(ctx context.Context, id, groupId int64) (m model.File, err error) {
	if groupId < 0 {
		return m, errcode.ArgsError().WithMsg("分组ID不合法")
	}

	db := s.DB(ctx, &m)
	err = db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return m, errcode.NotFound().WithMsg("文件不存在")
	}
	if err != nil {
		return m, err
	}

	if groupId > 0 {
		var count int64
		err = s.DB(ctx, &model.FileGroup{}).
			Model(&model.FileGroup{}).
			Where("id = ?", groupId).
			Count(&count).Error
		if err != nil {
			return m, err
		}
		if count == 0 {
			return m, errcode.ArgsError().WithMsg("文件分组不存在")
		}
	}

	err = db.Model(&m).Update("group_id", groupId).Error
	if err != nil {
		return m, err
	}
	err = db.Preload("Group").First(&m, id).Error
	return m, err
}

// fileObjectPath 生成文件存储路径
func fileObjectPath(originName string) string {
	ext := strings.ToLower(path.Ext(originName))
	return path.Join(time.Now().Format("2006/01/02"), uuid.New().String()+ext)
}
