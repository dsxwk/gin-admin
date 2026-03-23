package request

import (
	"errors"
	"github.com/gookit/validate"
)

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
	ID       int64  `json:"id" validate:"required|int|gt:0" label:"ID"`
	Username string `json:"username" validate:"required" label:"用户名"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Gender   int    `json:"gender" validate:"required|int" label:"性别"`
	Age      int    `json:"age" validate:"int" label:"年龄"`
}

// UserDetail 用户详情验证
type UserDetail struct {
	ID int64 `json:"id" validate:"required|int|gt:0" label:"ID"`
}

// UserDelete 用户删除验证
type UserDelete struct {
	ID int64 `json:"id" validate:"required|int|gt:0" label:"ID"`
}

// UserSearch 用户搜索
type UserSearch struct {
	Username string `form:"username" validate:"required" label:"用户名"`
	FullName string `form:"fullName" validate:"required" label:"姓名"`
	Nickname string `form:"nickname" validate:"required" label:"昵称"`
	Gender   int    `form:"gender" validate:"required|int" label:"性别"`
}

// User 用户请求验证
type User struct {
	ID       int64  `json:"id" validate:"required|int|gt:0" label:"ID"`
	Username string `json:"username" validate:"required" label:"用户名"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Gender   int    `json:"gender" validate:"required|int" label:"性别"`
	Password string `json:"password" validate:"required" label:"密码"`
	Age      int    `json:"age" validate:"int" label:"年龄"`
	PageListValidate
	Context
}

// UserFillAble 允许更新的键
var UserFillAble = []string{
	"avatar",
	"username",
	"fullName",
	"email",
	"password",
	"nickname",
	"gender",
	"age",
	"status",
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
// - 定义验证场景
// - 也可以添加验证设置
func (s User) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		// 列表
		"List": []string{
			"PageListValidate.Page",
			"PageListValidate.PageSize",
		},
		// 创建
		"Create": []string{
			"Username",
			"FullName",
			"Nickname",
			"Gender",
			"Password",
		},
		// 更新
		"Update": []string{
			"ID",
			"Username",
			"FullName",
			"Nickname",
			"Gender",
		},
		// 详情
		"Detail": []string{
			"ID",
		},
		// 删除
		"Delete": []string{
			"ID",
		},
	})
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
		"Page":     "页码",
		"PageSize": "每页数量",
		"ID":       "ID",
		"Username": "用户名",
		"FullName": "姓名",
		"Nickname": "昵称",
		"Gender":   "性别",
		"Password": "密码",
	}
}
