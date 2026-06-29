package request

import (
	"errors"
	"gin/common/base"
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
	ID       int64  `uri:"id" validate:"required|int|gt:0" label:"ID"`
	Username string `json:"username" validate:"required" label:"用户名"`
	FullName string `json:"fullName" validate:"required" label:"姓名"`
	Nickname string `json:"nickname" validate:"required" label:"昵称"`
	Gender   int    `json:"gender" validate:"required|int" label:"性别"`
	Age      int    `json:"age" validate:"int" label:"年龄"`
}

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
	PageListValidate
}

// UserRole 用户角色
type UserRole struct {
	UserId int64  `json:"userId" validate:"int" label:"用户id"`
	RoleId int64  `json:"roleId" validate:"int" label:"角色id"`
	Name   string `json:"name" validate:"" label:"角色名称"`
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
