package service

import (
	"fmt"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg"
	"time"
)

type DepartmentService struct {
	base.BaseService
}

// List 列表
func (s *DepartmentService) List(req request.Department) (pageData request.PageData, err error) {
	var (
		m      model.Department
		models []model.Department
		db     = s.DB(&m)
	)

	// 搜索
	db = s.Search(db, m, req.Search).
		Model(&m).
		Preload("DeptLeaders").
		Preload("DeptLeaders.Leader")

	err = db.Count(&pageData.Total).Error
	if err != nil {
		return pageData, err
	}

	if req.NotPage {
		err = db.Order("sort Asc").Find(&models).Error
		if err != nil {
			return pageData, err
		}
		pageData.List = m.GetTree(models)
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

// Create 创建
func (s *DepartmentService) Create(req request.Department) (model.Department, error) {
	var (
		m  model.Department
		db = s.DB(&m)
	)

	m = model.Department{
		Pid:    req.Pid,
		Name:   req.Name,
		Status: req.Status,
		Sort:   req.Sort,
	}

	err := db.Model(&m).Create(&m).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Update 更新
func (s *DepartmentService) Update(id int64, data map[string]any) (err error) {
	var (
		m  model.Department
		db = s.DB(&m)
	)

	tx := db.Begin()

	rows := model.FilterFields(db, m, data)
	rows[model.UpdatedField] = time.Now()
	if err = tx.Model(&m).Where("id = ?", id).Updates(rows).Error; err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.DepartmentLeaders{}).Where("department_id = ?", id).Delete(&model.DepartmentLeaders{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if pkg.HasKey(data, "deptLeaders") {
		deptLeaders, ok := data["deptLeaders"].([]any)
		if !ok {
			deptLeaders = []any{}
		}
		if len(deptLeaders) > 0 {
			var leaders []model.DepartmentLeaders
			for _, v := range deptLeaders {
				leaderMap, _ok := v.(map[string]any)
				if !_ok {
					continue
				}
				leaderUserId := int64(leaderMap["leaderUserId"].(float64))

				// 校验用户是否属于该部门
				var count int64
				tx.Model(&model.UserDepartments{}).
					Where("user_id = ? AND department_id = ?", leaderUserId, id).
					Count(&count)
				if count == 0 {
					tx.Rollback()
					return fmt.Errorf("用户%d不属于该部门,不能设为部门领导", leaderUserId)
				}

				leaders = append(leaders, model.DepartmentLeaders{
					DepartmentId: id,
					LeaderUserId: leaderUserId,
				})
			}
			if len(leaders) > 0 {
				err = tx.Model(&leaders).Create(&leaders).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	tx.Commit()

	return nil
}

// Detail 详情
func (s *DepartmentService) Detail(id int64) (m model.Department, err error) {
	var (
		db = s.DB(&m)
	)

	err = db.Model(&m).
		Preload("DeptLeaders").
		Preload("DeptLeaders.Leader").
		First(&m, id).Error
	if err != nil {
		return m, err
	}

	return m, nil
}

// Delete 删除
func (s *DepartmentService) Delete(id int64) (err error) {
	var (
		m          model.Department
		deptLeader model.DepartmentLeaders
		db         = s.DB(&m)
	)

	tx := db.Begin()

	err = tx.Model(&m).Delete(&m, id).Error
	if err != nil {
		return err
	}

	err = tx.Model(&deptLeader).Where("department_id = ?", id).Delete(&deptLeader).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
