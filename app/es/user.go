package es

import (
	"context"
	"errors"
	"fmt"
	"gin/app/facade"
	"gin/app/model"
	"gin/app/request"
	eslib "gin/pkg/serviceprovider/es"
	"strconv"
)

var userUpdateFields = map[string]bool{
	"avatar":    true,
	"username":  true,
	"fullName":  true,
	"email":     true,
	"nickname":  true,
	"gender":    true,
	"age":       true,
	"status":    true,
	"userRoles": true,
	"mainDept":  true,
	"userDepts": true,
	"createdAt": true,
	"updatedAt": true,
	"deletedAt": true,
}

var userTextFields = []string{"username", "fullName", "nickname", "email"}

// UserSearch 用户ES搜索服务
type UserSearch struct{}

// Client 获取ES客户端
func (s *UserSearch) Client() *eslib.Client {
	return facade.ES()
}

// IndexName 获取索引名称
func (s *UserSearch) IndexName() string {
	return model.TableNameUser
}

// Mapping 获取索引字段映射
func (s *UserSearch) Mapping() map[string]string {
	return map[string]string{
		"id":        "long",
		"avatar":    "keyword",
		"username":  "keyword",
		"fullName":  "keyword",
		"email":     "keyword",
		"nickname":  "text",
		"gender":    "byte",
		"age":       "integer",
		"status":    "byte",
		"userRoles": "nested",
		"mainDept":  "object",
		"userDepts": "nested",
		"createdAt": "date",
		"updatedAt": "date",
		"deletedAt": "date",
	}
}

// CreateIndex 创建用户索引
func (s *UserSearch) CreateIndex(ctx context.Context) error {
	return s.Client().CreateIndex(ctx, s.IndexName(), s.Mapping())
}

// EnsureIndex 确保用户索引存在
func (s *UserSearch) EnsureIndex(ctx context.Context) error {
	exists, err := s.Client().IndexExists(ctx, s.IndexName())
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.Client().CreateIndex(ctx, s.IndexName(), s.Mapping())
}

// DeleteIndex 删除用户索引
func (s *UserSearch) DeleteIndex(ctx context.Context) error {
	return s.Client().DeleteIndex(ctx, s.IndexName())
}

// Save 保存用户文档
func (s *UserSearch) Save(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("用户数据不能为空")
	}
	if user.ID <= 0 {
		return errors.New("用户ID必须大于0")
	}

	return s.Client().Index(ctx, s.IndexName(), strconv.FormatInt(user.ID, 10), s.Doc(user))
}

// Update 更新用户文档
func (s *UserSearch) Update(ctx context.Context, id int64, data map[string]any) error {
	if id <= 0 {
		return errors.New("用户ID必须大于0")
	}

	doc := make(map[string]any, len(data))
	for field, value := range data {
		if userUpdateFields[field] {
			doc[field] = value
		}
	}
	if len(doc) == 0 {
		return errors.New("没有可更新的ES字段")
	}

	return s.Client().Update(ctx, s.IndexName(), strconv.FormatInt(id, 10), doc)
}

// Delete 删除用户文档
func (s *UserSearch) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("用户ID必须大于0")
	}
	return s.Client().Delete(ctx, s.IndexName(), strconv.FormatInt(id, 10))
}

// Detail 获取用户文档
func (s *UserSearch) Detail(ctx context.Context, id int64) (map[string]any, error) {
	if id <= 0 {
		return nil, errors.New("用户ID必须大于0")
	}

	doc, err := s.Client().GetDocument[map[string]any](ctx, s.IndexName(), strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	if doc == nil || !doc.Found {
		return nil, fmt.Errorf("用户ES文档不存在,ID:%d", id)
	}
	return doc.Source, nil
}

// List 分页搜索用户
func (s *UserSearch) List(ctx context.Context, conditions map[string]any, page, pageSize int, sorts map[string]any) (request.PageData, error) {
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

	result, err := s.Client().Search[map[string]any](ctx, s.IndexName(), body)
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
func (s *UserSearch) Query(conditions map[string]any) map[string]any {
	must := make([]any, 0, len(conditions))
	for field, value := range conditions {
		if field == "keyword" {
			must = append(must, map[string]any{
				"multi_match": map[string]any{
					"query":  value,
					"fields": userTextFields,
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

// Doc 将用户模型转换为ES文档
func (s *UserSearch) Doc(user *model.User) map[string]any {
	if user == nil {
		return nil
	}

	userRoles := make([]map[string]any, 0, len(user.UserRoles))
	for _, role := range user.UserRoles {
		if role == nil {
			continue
		}
		userRoles = append(userRoles, map[string]any{
			"id":     role.ID,
			"userId": role.UserID,
			"roleId": role.RoleID,
			"name":   role.Name,
		})
	}

	var mainDept map[string]any
	if user.MainDept != nil {
		mainDept = departmentDoc(user.MainDept)
	}

	userDepts := make([]map[string]any, 0, len(user.UserDepts))
	for _, department := range user.UserDepts {
		if department == nil {
			continue
		}
		userDepts = append(userDepts, departmentDoc(department))
	}

	return map[string]any{
		"id":        user.ID,
		"avatar":    user.Avatar,
		"username":  user.Username,
		"fullName":  user.FullName,
		"email":     user.Email,
		"nickname":  user.Nickname,
		"gender":    user.Gender,
		"age":       user.Age,
		"status":    user.Status,
		"userRoles": userRoles,
		"mainDept":  mainDept,
		"userDepts": userDepts,
		"createdAt": s.Client().FormatDateTime(user.CreatedAt),
		"updatedAt": s.Client().FormatDateTime(user.UpdatedAt),
	}
}

// departmentDoc 部门文档
func departmentDoc(department *model.UserDepartments) map[string]any {
	doc := map[string]any{
		"id":           department.ID,
		"userId":       department.UserId,
		"departmentId": department.DepartmentId,
		"isMain":       department.IsMain,
	}
	if department.Dept != nil {
		doc["name"] = department.Dept.Name
	}
	return doc
}
