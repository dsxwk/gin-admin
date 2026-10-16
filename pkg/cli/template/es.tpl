package {{.Package}}

import (
	"context"
	"errors"
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/app/request"
	"strconv"
)

var {{.Var}}UpdateFields = map[string]bool{
{{.UpdateFields}}
}

var {{.Var}}TextFields = []string{
{{- if .TextFields}}
{{.TextFields}}
{{- end}}
}

// {{.Name}}Search {{.Description}}ES搜索服务
type {{.Name}}Search struct {}

// IndexName 获取索引名称
func (s *{{.Name}}Search) IndexName() string {
	return "{{.Table}}"
}

// Mapping 获取索引字段映射
func (s *{{.Name}}Search) Mapping() map[string]string {
	return map[string]string{
{{.MappingFields}}
	}
}

// CreateIndex 创建{{.Description}}索引
func (s *{{.Name}}Search) CreateIndex(ctx context.Context) error {
	return facade.ES().CreateIndex(ctx, s.IndexName(), s.Mapping())
}

// EnsureIndex 确保{{.Description}}索引存在
func (s *{{.Name}}Search) EnsureIndex(ctx context.Context) error {
	exists, err := facade.ES().IndexExists(ctx, s.IndexName())
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return facade.ES().CreateIndex(ctx, s.IndexName(), s.Mapping())
}

// DeleteIndex 删除{{.Description}}索引
func (s *{{.Name}}Search) DeleteIndex(ctx context.Context) error {
	return facade.ES().DeleteIndex(ctx, s.IndexName())
}

// Save 保存{{.Description}}文档
func (s *{{.Name}}Search) Save(ctx context.Context, m *model.{{.Name}}) error {
	if m == nil {
		return errors.New("{{.Description}}数据不能为空")
	}
	if m.ID <= 0 {
		return errors.New("{{.Description}}ID必须大于0")
	}

	return facade.ES().Index(ctx, s.IndexName(), strconv.FormatInt(m.ID, 10), s.Doc(m))
}

// Update 更新{{.Description}}文档
func (s *{{.Name}}Search) Update(ctx context.Context, id int64, data map[string]any) error {
	if id <= 0 {
		return errors.New("{{.Description}}ID必须大于0")
	}

	doc := make(map[string]any, len(data))
	for field, value := range data {
		if {{.Var}}UpdateFields[field] {
			doc[field] = value
		}
	}
	if len(doc) == 0 {
		return errors.New("没有可更新的ES字段")
	}

	return facade.ES().Update(ctx, s.IndexName(), strconv.FormatInt(id, 10), doc)
}

// Delete 删除{{.Description}}文档
func (s *{{.Name}}Search) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("{{.Description}}ID必须大于0")
	}
	return facade.ES().Delete(ctx, s.IndexName(), strconv.FormatInt(id, 10))
}

// Detail 获取{{.Description}}文档
func (s *{{.Name}}Search) Detail(ctx context.Context, id int64) (map[string]any, error) {
	if id <= 0 {
		return nil, errors.New("{{.Description}}ID必须大于0")
	}

	doc, err := facade.ES().Document[map[string]any](ctx, s.IndexName(), strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	if doc == nil || !doc.Found {
		return nil, fmt.Errorf("{{.Description}}ES文档不存在,ID:%d", id)
	}
	return doc.Source, nil
}

// List 分页搜索{{.Description}}
func (s *{{.Name}}Search) List(ctx context.Context, conditions map[string]any, page, pageSize int, sorts map[string]any) (request.PageData, error) {
	offset, size := request.Pagination(page, pageSize)
	body := map[string]any{
		"from": offset,
		"size": size,
	}
	if len(conditions) > 0 {
		body["query"] = s.Query(conditions)
	}
	if len(sorts) > 0 {
		sortList := make([]any, 0, len(sorts))
		for field, direction := range sorts {
			sortList = append(sortList, map[string]any{field: direction})
		}
		body["sort"] = sortList
	}

	result, err := facade.ES().Search[map[string]any](ctx, s.IndexName(), body)
	if err != nil {
		return request.PageData{}, err
	}

	list := make([]map[string]any, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		if hit.Source != nil {
			list = append(list, hit.Source)
		}
	}

	return request.PageData{
		Total:    result.Hits.Total.Value,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// Query 将简单条件转换为ES查询体
func (s *{{.Name}}Search) Query(conditions map[string]any) map[string]any {
	must := make([]any, 0, len(conditions))
	for field, value := range conditions {
		if field == "keyword" {
			must = append(must, map[string]any{
				"multi_match": map[string]any{
					"query":  value,
					"fields": {{.Var}}TextFields,
				},
			})
			continue
		}

		if values, ok := value.([]any); ok {
			must = append(must, map[string]any{
				"terms": map[string]any{field: values},
			})
			continue
		}

		must = append(must, map[string]any{
			"term": map[string]any{field: value},
		})
	}

	if len(must) == 1 {
		if query, ok := must[0].(map[string]any); ok {
			return query
		}
	}
	return map[string]any{
		"bool": map[string]any{
			"must": must,
		},
	}
}

// Doc 将{{.Description}}模型转换为ES文档
func (s *{{.Name}}Search) Doc(m *model.{{.Name}}) map[string]any {
	if m == nil {
		return nil
	}

	return map[string]any{
{{.DocFields}}
	}
}
