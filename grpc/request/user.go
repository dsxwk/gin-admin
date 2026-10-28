package request

import "github.com/dsxwk/validatex"

// UserRequest 用户请求
type UserRequest struct {
	Id int32 `json:"id" label:"ID"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int            `json:"page" label:"页码"`
	PageSize int            `json:"pageSize" label:"每页数量"`
	NotPage  bool           `json:"notPage" label:"不分页"`
	Search   map[string]any `json:"search" label:"搜索条件"`
	Sort     map[string]any `json:"sort" label:"排序条件"`
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username  string         `json:"username" label:"用户名"`
	FullName  string         `json:"fullName" label:"姓名"`
	Nickname  string         `json:"nickname" label:"昵称"`
	Gender    int32          `json:"gender" label:"性别"`
	Password  string         `json:"password" label:"密码"`
	Age       int32          `json:"age" label:"年龄"`
	UserRoles []UserRoleItem `json:"userRoles" label:"用户角色"`
	MainDept  DeptItem       `json:"mainDept" label:"主部门"`
	UserDepts []DeptItem     `json:"userDepts" label:"用户部门"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Id        int32          `json:"id" label:"ID"`
	Username  string         `json:"username" label:"用户名"`
	FullName  string         `json:"fullName" label:"姓名"`
	Nickname  string         `json:"nickname" label:"昵称"`
	Gender    int32          `json:"gender" label:"性别"`
	Password  string         `json:"password" label:"密码"`
	Age       int32          `json:"age" label:"年龄"`
	UserRoles []UserRoleItem `json:"userRoles" label:"用户角色"`
	MainDept  DeptItem       `json:"mainDept" label:"主部门"`
	UserDepts []DeptItem     `json:"userDepts" label:"用户部门"`
}

// UserBatchDeleteRequest 用户批量删除请求
type UserBatchDeleteRequest struct {
	Ids []int32 `json:"ids" label:"ID列表"`
}

// UserRoleItem 用户角色项
type UserRoleItem struct {
	RoleId int32  `json:"roleId" label:"角色ID"`
	Name   string `json:"name" label:"角色名称"`
}

// DeptItem 部门项
type DeptItem struct {
	DepartmentId int32 `json:"departmentId" label:"部门ID"`
}

// Validate UserRoleItem验证
func (s UserRoleItem) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Int("roleId", s.RoleId).Required().Gt(0).Errors()
}

// Validate DeptItem验证
func (s DeptItem) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Int("departmentId", s.DepartmentId).Required().Gt(0).Errors()
}

// ValidateWithScene 用户请求验证
func (s UserRequest) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Detail", "Delete":
		return fields.Int("id", s.Id).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Validate 用户列表请求验证
func (s UserListRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.When(!s.NotPage, func() validatex.Errors {
		return validatex.Check(
			fields.Int("page", s.Page).Required().Gt(0).Errors(),
			fields.Int("pageSize", s.PageSize).Required().Gt(0).Errors(),
		)
	})
}

// Validate 用户创建请求验证
func (s UserCreateRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Str("username", s.Username).Required().Errors(),
		fields.Str("fullName", s.FullName).Required().Errors(),
		fields.Str("nickname", s.Nickname).Required().Errors(),
		fields.Int("gender", s.Gender).Required().Errors(),
		fields.Str("password", s.Password).Required().Errors(),
		fields.Struct("mainDept", s.MainDept),
		fields.Slice[struct{}]("userDepts", make([]struct{}, len(s.UserDepts))).Required().Errors(),
		fields.SliceStruct("userDepts", s.UserDepts),
	)
}

// Validate 用户更新请求验证
func (s UserUpdateRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("id", s.Id).Required().Gt(0).Errors(),
		fields.Str("fullName", s.FullName).Required().Errors(),
	)
}

// Validate 用户批量删除请求验证
func (s UserBatchDeleteRequest) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Slice("ids", s.Ids).Required().Each("gt", "", func(id int32) bool {
		return id > 0
	}, int32(0)).Errors()
}

// userMessages 用户请求公共错误消息
func userMessages() validatex.Messages {
	return validatex.Messages{
		"required": "字段 {field} 必填",
		"gt":       "字段 {field} 需大于 0",
	}
}

// userTranslates 用户请求公共字段翻译
func userTranslates() validatex.Translates {
	return validatex.Translates{
		"id":                       "ID",
		"page":                     "页码",
		"pageSize":                 "每页数量",
		"username":                 "用户名",
		"fullName":                 "姓名",
		"nickname":                 "昵称",
		"gender":                   "性别",
		"password":                 "密码",
		"age":                      "年龄",
		"mainDept.departmentId":    "主部门",
		"userDepts":                "用户部门",
		"userDepts.*.departmentId": "部门ID",
		"ids":                      "ID列表",
		"ids.*":                    "ID",
		"roleId":                   "角色ID",
		"departmentId":             "部门ID",
	}
}

// Messages 用户请求错误消息
func (s UserRequest) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户请求字段翻译
func (s UserRequest) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户列表错误消息
func (s UserListRequest) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户列表字段翻译
func (s UserListRequest) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户创建错误消息
func (s UserCreateRequest) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户创建字段翻译
func (s UserCreateRequest) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户更新错误消息
func (s UserUpdateRequest) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户更新字段翻译
func (s UserUpdateRequest) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户批量删除错误消息
func (s UserBatchDeleteRequest) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户批量删除字段翻译
func (s UserBatchDeleteRequest) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户角色项错误消息
func (s UserRoleItem) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户角色项字段翻译
func (s UserRoleItem) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 部门项错误消息
func (s DeptItem) Messages() validatex.Messages {
	return userMessages()
}

// Translates 部门项字段翻译
func (s DeptItem) Translates() validatex.Translates {
	return userTranslates()
}
