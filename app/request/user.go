package request

import (
	"errors"
	"gin/common/base"
	"gin/pkg"
	"github.com/gookit/validate"
)

// User 用户请求验证
type User struct {
	base.BaseRequest
	ID        int64      `uri:"id" validate:"required|int|gt:0" label:"ID"`
	Username  string     `json:"username" form:"username" validate:"required" label:"用户名"`
	FullName  string     `json:"fullName" validate:"required" label:"姓名"`
	Nickname  string     `json:"nickname" validate:"required" label:"昵称"`
	Gender    int64      `json:"gender" validate:"required|int" label:"性别"`
	Password  string     `json:"password" validate:"required" label:"密码"`
	Age       int64      `json:"age" validate:"int" label:"年龄"`
	UserRoles []UserRole `json:"userRoles" validate:"" label:"用户角色"`
	IDs       []int64    `json:"ids" validate:"required" label:"ID列表"`
	PageListValidate
}

// UserCreate 用户创建验证
type UserCreate struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Gender   int    `json:"gender" validate:"required|int" label:"性别"`
	Password string `json:"password" validate:"required" label:"密码"`
}

// UserUpdate 用户更新验证
type UserUpdate struct {
	ID       int64  `uri:"id" validate:"required|int|gt:0" label:"ID"`
	Username string `json:"username" validate:"required" label:"用户名"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Gender   int    `json:"gender" validate:"required|int" label:"性别"`
	Age      int    `json:"age" validate:"int" label:"年龄"`
}

// UserPassword 用户密码
type UserPassword struct {
	Password string `json:"password" validate:"required" label:"密码"`
}

// UserRole 用户角色
type UserRole struct {
	UserId int64  `json:"userId" validate:"int" label:"用户id"`
	RoleId int64  `json:"roleId" validate:"int" label:"角色id"`
	Name   string `json:"name" validate:"" label:"角色名称"`
}

// UserImport 用户导入
type UserImport struct {
	Data []UserImportItem `json:"data" validate:"required|minLen:1" label:"用户列表"`
}

// UserImportItem 用户导入子项
type UserImportItem struct {
	Username string `json:"username" validate:"required|minLen:3|maxLen:20|regex:^[a-zA-Z0-9_]+$" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Email    string `json:"email" validate:"required|email" label:"邮箱"`
	Gender   int64  `json:"gender" validate:"required|int" label:"性别"`
	Age      int64  `json:"age" validate:"int" label:"年龄"`
	Status   int64  `json:"status" validate:"int" label:"状态"`
}

// UserBatchDelete 用户批量删除
type UserBatchDelete struct {
	IDs []int64 `json:"ids" validate:"required" label:"ID列表"`
}

// Validate 请求验证
func (s User) Validate(data User, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s User) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":        []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create":      []string{"Username", "FullName", "Nickname", "Gender", "Password"},
		"Update":      []string{"ID", "Username", "FullName", "Nickname", "Gender"},
		"Detail":      []string{"ID"},
		"Delete":      []string{"ID"},
		"BatchDelete": []string{"IDs"},
		"Password":    []string{"ID", "Password"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s User) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"PageListValidate.Page.gt":     "字段 {field} 需大于 0",
		"PageListValidate.PageSize.gt": "字段 {field} 需大于 0",
	}
}

// Translates 字段翻译
func (s User) Translates() map[string]string {
	return validate.MS{
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
		"ID":                        "ID",
		"Username":                  "用户名",
		"FullName":                  "姓名",
		"Nickname":                  "昵称",
		"Gender":                    "性别",
		"Password":                  "密码",
	}
}

// Validate 用户导入验证
func (s UserImport) Validate(data UserImport, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s UserImport) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"Import": []string{"Data"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s UserImport) Messages() map[string]string {
	return validate.MS{
		"required": "{field} 必填",
		"minLen":   "{field} 长度不能少于 {min} 个字符",
		"maxLen":   "{field} 长度不能超过 {max} 个字符",
		"int":      "{field} 必须为整数",
		"regex":    "{field} 格式错误",
		"email":    "{field} 邮箱格式错误",
	}
}

// Translates 字段翻译
func (s UserImport) Translates() map[string]string {
	ms := validate.MS{
		"Data": "导入数据",
	}
	for i := range s.Data {
		prefix := pkg.Sprintf("Data.%d.", i)
		rowLabel := pkg.Sprintf("第 %d 行 ", i+1)
		ms[prefix+"Username"] = rowLabel + "用户名"
		ms[prefix+"Password"] = rowLabel + "密码"
		ms[prefix+"FullName"] = rowLabel + "姓名"
		ms[prefix+"Nickname"] = rowLabel + "昵称"
		ms[prefix+"Email"] = rowLabel + "邮箱"
		ms[prefix+"Gender"] = rowLabel + "性别"
		ms[prefix+"Age"] = rowLabel + "年龄"
		ms[prefix+"Status"] = rowLabel + "状态"
	}
	return ms
}
