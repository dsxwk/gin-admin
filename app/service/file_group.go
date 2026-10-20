package service

import (
	"context"
	"errors"
	"gin/app/errcode"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"

	"gorm.io/gorm"
)

// FileGroupService 文件分组服务
type FileGroupService struct {
	base.Service
}

// List 文件分组列表
func (s *FileGroupService) List(ctx context.Context, req request.FileGroup) (pageData request.PageData, err error) {
	var (
		m      model.FileGroup
		models []model.FileGroup
		db     = s.DB(ctx, &m)
	)

	db = s.Search(db, m, req.Search).Model(&m)
	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("sort ASC").Order("id ASC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	} else {
		pageData.Page = req.Page
		pageData.PageSize = req.PageSize
		offset, limit := request.Pagination(req.Page, req.PageSize)

		err = db.Offset(offset).Limit(limit).Order("sort ASC").Order("id ASC").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = models
	}

	return pageData, nil
}

// Create 创建文件分组
func (s *FileGroupService) Create(ctx context.Context, req request.FileGroup) (m model.FileGroup, err error) {
	if err = s.validateParent(ctx, 0, req.Pid); err != nil {
		return m, err
	}

	m = model.FileGroup{
		Pid:  req.Pid,
		Name: req.Name,
		Sort: req.Sort,
	}
	if err = s.DB(ctx, &m).Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// Update 修改文件分组
func (s *FileGroupService) Update(ctx context.Context, id int64, req request.FileGroup, data map[string]any) error {
	if err := s.validateParent(ctx, id, req.Pid); err != nil {
		return err
	}
	return s.UpdateByID(ctx, &model.FileGroup{}, id, data)
}

// Delete 删除文件分组
func (s *FileGroupService) Delete(ctx context.Context, id int64) error {
	var m model.FileGroup
	db := s.DB(ctx, &m)
	err := db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errcode.NotFound().WithMsg("文件分组不存在")
	}
	if err != nil {
		return err
	}

	var childCount int64
	err = db.Model(&model.FileGroup{}).Where("pid = ?", id).Count(&childCount).Error
	if err != nil {
		return err
	}
	if childCount > 0 {
		return errcode.ArgsError().WithMsg("请先删除子分组")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.File{}).Where("group_id = ?", id).Update("group_id", 0).Error
		if err != nil {
			return err
		}
		return tx.Delete(&model.FileGroup{}, id).Error
	})
}

// validateParent 校验父级分组
func (s *FileGroupService) validateParent(ctx context.Context, id, parentId int64) error {
	for parentId > 0 {
		if id > 0 && parentId == id {
			return errcode.ArgsError().WithMsg("父级分组不能是当前分组")
		}

		var parent model.FileGroup
		err := s.DB(ctx, &parent).Select("id", "pid").First(&parent, parentId).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ArgsError().WithMsg("父级分组不存在")
		}
		if err != nil {
			return err
		}
		parentId = parent.Pid
	}
	return nil
}
