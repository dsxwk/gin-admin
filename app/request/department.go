package request

import (
	"gin/common/base"
	"gin/common/errcode"

	"github.com/gookit/validate"
)

// Department 请求验证
type Department struct {
	base.BaseRequest
	ID          int64        `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Pid         int64        `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name        string       `json:"name" form:"name" validate:"required" label:"部门名称"`
	Status      int64        `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort        int64        `json:"sort" form:"sort" validate:"int" label:"排序"`
	DeptLeaders []DeptLeader `json:"deptLeaders" form:"deptLeaders" validate:"required" label:"部门领导"`
	PageListValidate
}

// DeptLeader 部门领导
type DeptLeader struct {
	DepartmentId int64 `json:"departmentId" form:"departmentId" validate:"required|int|gt:0" label:"部门ID"`
	LeaderUserId int64 `json:"leaderUserId" form:"leaderUserId" validate:"required|int|gt:0" label:"部门领导ID"`
}

// Validate 请求验证
func (s Department) Validate(data Department, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// DepartmentCreate Department创建验证
type DepartmentCreate struct {
	Pid    int64  `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name   string `json:"name" form:"name" validate:"required" label:"部门名称"`
	Status int64  `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" validate:"int" label:"排序"`
}

// DepartmentUpdate Department更新验证
type DepartmentUpdate struct {
	Pid    int64  `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name   string `json:"name" form:"name" validate:"required" label:"部门名称"`
	Status int64  `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" validate:"int" label:"排序"`
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Department) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List": []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{
			"Pid",
			"Name",
			"Status",
			"Sort",
		},
		"Update": []string{
			"ID",
			"Pid",
			"Name",
			"Status",
			"Sort",
		},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Department) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"minLen":                       "{field} 长度不能少于 {min} 个字符",
		"maxLen":                       "{field} 长度不能超过 {max} 个字符",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s Department) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Pid":                       "父级id",
		"Name":                      "部门名称",
		"Status":                    "状态 1=启用 2=停用",
		"Sort":                      "排序",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
