package {{.Package}}

import (
    "gin/app/model"
    "gin/app/request"
    "gin/common/base"
    "time"
)

type {{.Name}}Service struct {
    base.BaseService
}

// List 列表
func (s *{{.Name}}Service) List(req request.{{.Name}}) (pageData request.PageData, err error) {
    var (
        m      model.{{.Name}}
        models []model.{{.Name}}
        db     = s.DB(&m)
    )

    // 搜索
    db = s.Search(db, m, req.Search).Model(&m)

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

// Create 创建
func (s *{{.Name}}Service) Create(req request.{{.Name}}) (request.{{.Name}}, error) {
    var (
        m  model.{{.Name}}
        db = s.DB(&m)
    )
    {{if .HasFields}}
    m = model.{{.Name}}{
{{.Fields}}
    }
    {{else}}
    // todo
    {{end}}
    err := db.Model(&m).Create(&m).Error
    if err != nil {
        return req, err
    }

    return req, nil
}

// Update 更新
func (s *{{.Name}}Service) Update(id int64, data map[string]interface{}) (err error) {
    var (
        m  model.{{.Name}}
        db = s.DB(&m)
    )

    rows := model.FilterFields(db, m, data)
    rows["updated_at"] = time.DateTime

    err = db.Model(&m).Where("id = ?", id).Updates(rows).Error
    if err != nil {
        return err
    }

    return nil
}

// Detail 详情
func (s *{{.Name}}Service) Detail(id int64) (m model.{{.Name}}, err error) {
    var (
        db = s.DB(&m)
    )

    err = db.Model(&m).First(&m, id).Error
    if err != nil {
        return m, err
    }

    return m, nil
}

// Delete 删除
func (s *{{.Name}}Service) Delete(id int64) (err error) {
    var (
        m  model.{{.Name}}
        db = s.DB(&m)
    )

    err = db.Model(&m).Delete(&m, id).Error
    if err != nil {
        return err
    }

    return nil
}
