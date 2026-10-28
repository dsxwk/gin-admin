package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Department 请求验证
type Department struct {
	base.Request
	ID          int64        `json:"id" form:"id" label:"ID"`
	Pid         int64        `json:"pid" form:"pid" label:"父级id"`
	Name        string       `json:"name" form:"name" label:"部门名称"`
	Status      int64        `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort        int64        `json:"sort" form:"sort" label:"排序"`
	DeptLeaders []DeptLeader `json:"deptLeaders" form:"deptLeaders" label:"部门领导"`
	PageListValidate
}

// DeptLeader 部门领导
type DeptLeader struct {
	DepartmentId int64 `json:"departmentId" form:"departmentId" label:"部门ID"`
	LeaderUserId int64 `json:"leaderUserId" form:"leaderUserId" label:"部门领导ID"`
}

// DepartmentCreate Department创建验证
type DepartmentCreate struct {
	Pid    int64  `json:"pid" form:"pid" label:"父级id"`
	Name   string `json:"name" form:"name" label:"部门名称"`
	Status int64  `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" label:"排序"`
}

// DepartmentUpdate Department更新验证
type DepartmentUpdate struct {
	Pid    int64  `json:"pid" form:"pid" label:"父级id"`
	Name   string `json:"name" form:"name" label:"部门名称"`
	Status int64  `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" label:"排序"`
}

// Validate 部门领导验证
func (s DeptLeader) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("departmentId", s.DepartmentId).Required().Gt(0).Errors(),
		fields.Int("leaderUserId", s.LeaderUserId).Required().Gt(0).Errors(),
	)
}

// ValidateWithScene 请求验证
func (s Department) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Str("name", s.Name).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Department) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Department) Translates() validatex.Translates {
	return validatex.Translates{
		"id":           "ID",
		"pid":          "父级id",
		"name":         "部门名称",
		"status":       "状态 1=启用 2=停用",
		"sort":         "排序",
		"departmentId": "部门ID",
		"leaderUserId": "部门领导ID",
		"page":         "页码",
		"pageSize":     "每页数量",
	}
}
