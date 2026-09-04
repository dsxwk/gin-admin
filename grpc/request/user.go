package request

import (
	"gin/common/errcode"

	"github.com/gookit/validate"
)

// UserRequest 用户请求
type UserRequest struct {
	Id int32 `json:"id" validate:"required|gt:0" label:"ID"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int            `json:"page" validate:"required|int|gt:0" label:"页码"`
	PageSize int            `json:"pageSize" validate:"required|int|gt:0" label:"每页数量"`
	NotPage  bool           `json:"notPage" label:"不分页"`
	Search   map[string]any `json:"search" label:"搜索条件"`
	Sort     map[string]any `json:"sort" label:"排序条件"`
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username  string         `json:"username" validate:"required" label:"用户名"`
	FullName  string         `json:"fullName" validate:"required" label:"姓名"`
	Nickname  string         `json:"nickname" validate:"required" label:"昵称"`
	Gender    int32          `json:"gender" validate:"required|int" label:"性别"`
	Password  string         `json:"password" validate:"required" label:"密码"`
	Age       int32          `json:"age" validate:"int" label:"年龄"`
	UserRoles []UserRoleItem `json:"userRoles" label:"用户角色"`
	MainDept  DeptItem       `json:"mainDept" validate:"required" label:"主部门"`
	UserDepts []DeptItem     `json:"userDepts" validate:"required" label:"用户部门"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Id        int32          `json:"id" validate:"required|int|gt:0" label:"ID"`
	Username  string         `json:"username" validate:"" label:"用户名"`
	FullName  string         `json:"fullName" validate:"required" label:"姓名"`
	Nickname  string         `json:"nickname" validate:"" label:"昵称"`
	Gender    int32          `json:"gender" validate:"int" label:"性别"`
	Password  string         `json:"password" label:"密码"`
	Age       int32          `json:"age" validate:"int" label:"年龄"`
	UserRoles []UserRoleItem `json:"userRoles" label:"用户角色"`
	MainDept  DeptItem       `json:"mainDept" label:"主部门"`
	UserDepts []DeptItem     `json:"userDepts" label:"用户部门"`
}

// UserBatchDeleteRequest 用户批量删除请求
type UserBatchDeleteRequest struct {
	Ids []int32 `json:"ids" validate:"required|minLen:1" label:"ID列表"`
}

// UserRoleItem 用户角色项
type UserRoleItem struct {
	RoleId int32  `json:"roleId" validate:"required|int|gt:0" label:"角色ID"`
	Name   string `json:"name" label:"角色名称"`
}

// DeptItem 部门项
type DeptItem struct {
	DepartmentId int32 `json:"departmentId" validate:"required|int|gt:0" label:"部门ID"`
}

// Validate 请求验证
func (s UserRequest) Validate(data UserRequest, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s UserRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Detail": []string{"Id"},
		"Delete": []string{"Id"},
	})
}

// Messages 验证器错误消息
func (s UserRequest) Messages() map[string]string {
	return userMessages()
}

// Translates 字段翻译
func (s UserRequest) Translates() map[string]string {
	return userTranslates()
}

// Validate 列表请求验证
func (s UserListRequest) Validate() error {
	return validateRequest(s, "List")
}

// ConfigValidation 配置验证
func (s UserListRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"List": []string{"Page", "PageSize"},
	})
}

// Messages 验证器错误消息
func (s UserListRequest) Messages() map[string]string {
	return userMessages()
}

// Translates 字段翻译
func (s UserListRequest) Translates() map[string]string {
	return userTranslates()
}

// Validate 创建请求验证
func (s UserCreateRequest) Validate() error {
	return validateRequest(s, "Create")
}

// ConfigValidation 配置验证
func (s UserCreateRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Create": []string{"Username", "FullName", "Nickname", "Gender", "Password", "MainDept.DepartmentId", "UserDepts"},
	})
}

// Messages 验证器错误消息
func (s UserCreateRequest) Messages() map[string]string {
	return userMessages()
}

// Translates 字段翻译
func (s UserCreateRequest) Translates() map[string]string {
	return userTranslates()
}

// Validate 更新请求验证
func (s UserUpdateRequest) Validate() error {
	return validateRequest(s, "Update")
}

// ConfigValidation 配置验证
func (s UserUpdateRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Update": []string{"Id", "Username", "FullName", "Nickname", "Gender"},
	})
}

// Messages 验证器错误消息
func (s UserUpdateRequest) Messages() map[string]string {
	return userMessages()
}

// Translates 字段翻译
func (s UserUpdateRequest) Translates() map[string]string {
	return userTranslates()
}

// Validate 批量删除请求验证
func (s UserBatchDeleteRequest) Validate() error {
	return validateRequest(s, "BatchDelete")
}

// ConfigValidation 配置验证
func (s UserBatchDeleteRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"BatchDelete": []string{"Ids"},
	})
}

// Messages 验证器错误消息
func (s UserBatchDeleteRequest) Messages() map[string]string {
	return userMessages()
}

// Translates 字段翻译
func (s UserBatchDeleteRequest) Translates() map[string]string {
	return userTranslates()
}

// validateRequest 通用请求验证
func validateRequest(data any, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// userMessages 用户请求公共错误消息
func userMessages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"gt":       "字段 {field} 需大于 0",
		"minLen":   "字段 {field} 不能为空",
	}
}

// userTranslates 用户请求公共字段翻译
func userTranslates() map[string]string {
	return validate.MS{
		"Id":                    "ID",
		"Page":                  "页码",
		"PageSize":              "每页数量",
		"Username":              "用户名",
		"FullName":              "姓名",
		"Nickname":              "昵称",
		"Gender":                "性别",
		"Password":              "密码",
		"Age":                   "年龄",
		"MainDept.DepartmentId": "主部门",
		"UserDepts":             "用户部门",
		"Ids":                   "ID列表",
	}
}
