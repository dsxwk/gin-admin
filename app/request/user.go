package request

import (
	"gin/common/base"
	"regexp"

	"github.com/dsxwk/validatex"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// User 用户请求验证
type User struct {
	base.Request
	ID        int64            `uri:"id" form:"id" label:"ID"`
	Username  string           `json:"username" form:"username" label:"用户名"`
	FullName  string           `json:"fullName" form:"fullName" label:"姓名"`
	Nickname  string           `json:"nickname" form:"nickname" label:"昵称"`
	Gender    int64            `json:"gender" form:"gender" label:"性别"`
	Password  string           `json:"password" form:"password" label:"密码"`
	Age       int64            `json:"age" form:"age" label:"年龄"`
	UserRoles []UserRole       `json:"userRoles" form:"userRoles" label:"用户角色"`
	IDs       []int64          `json:"ids" form:"ids" label:"ID列表"`
	Data      []UserImportItem `json:"data" form:"data" label:"用户列表"`
	MainDept  Dept             `json:"mainDept" form:"mainDept" label:"主部门"`
	UserDepts []Dept           `json:"userDepts" form:"userDepts" label:"用户部门"`
	PageListValidate
}

// Dept 部门
type Dept struct {
	DepartmentId int64 `json:"departmentId" form:"departmentId" label:"部门id"`
}

// UserCreate 用户创建验证
type UserCreate struct {
	Username  string `json:"username" form:"username" label:"用户名"`
	FullName  string `json:"fullName" form:"fullName" label:"姓名"`
	Nickname  string `json:"nickname" form:"nickname" label:"昵称"`
	Gender    int    `json:"gender" form:"gender" label:"性别"`
	Password  string `json:"password" form:"password" label:"密码"`
	MainDept  Dept   `json:"mainDept" form:"mainDept" label:"主部门"`
	UserDepts []Dept `json:"userDepts" form:"userDepts" label:"用户部门"`
}

// UserUpdate 用户更新验证
type UserUpdate struct {
	Username  string `json:"username" form:"username" label:"用户名"`
	FullName  string `json:"fullName" form:"fullName" label:"姓名"`
	Nickname  string `json:"nickname" form:"nickname" label:"昵称"`
	Gender    int    `json:"gender" form:"gender" label:"性别"`
	Age       int    `json:"age" form:"age" label:"年龄"`
	MainDept  Dept   `json:"mainDept" form:"mainDept" label:"主部门"`
	UserDepts []Dept `json:"userDepts" form:"userDepts" label:"用户部门"`
}

// UserPassword 用户密码
type UserPassword struct {
	Password string `json:"password" form:"password" label:"密码"`
}

// UserRole 用户角色
type UserRole struct {
	UserId int64  `json:"userId" form:"userId" label:"用户id"`
	RoleId int64  `json:"roleId" form:"roleId" label:"角色id"`
	Name   string `json:"name" form:"name" label:"角色名称"`
}

// UserImport 用户导入
type UserImport struct {
	Data []UserImportItem `json:"data" form:"data" label:"用户列表"`
}

// UserImportItem 用户导入子项
type UserImportItem struct {
	Username string `json:"username" form:"username" label:"用户名"`
	Password string `json:"password" form:"password" label:"密码"`
	FullName string `json:"fullName" form:"fullName" label:"姓名"`
	Nickname string `json:"nickname" form:"nickname" label:"昵称"`
	Email    string `json:"email" form:"email" label:"邮箱"`
	Gender   int64  `json:"gender" form:"gender" label:"性别"`
	Age      int64  `json:"age" form:"age" label:"年龄"`
	Status   int64  `json:"status" form:"status" label:"状态"`
}

// UserBatchDelete 用户批量删除
type UserBatchDelete struct {
	IDs []int64 `json:"ids" form:"ids" label:"ID列表"`
}

// ValidateWithScene 用户请求验证
func (s User) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return validatex.When(!s.NotPage, func() validatex.Errors {
			return validatex.Check(
				fields.Int("PageListValidate.Page", s.Page).Required().Gt(0).Errors(),
				fields.Int("PageListValidate.PageSize", s.PageSize).Required().Gt(0).Errors(),
			)
		})
	case "Create":
		return validatex.Check(
			fields.Str("Username", s.Username).Required().Errors(),
			fields.Str("FullName", s.FullName).Required().Errors(),
			fields.Str("Nickname", s.Nickname).Required().Errors(),
			fields.Int("Gender", s.Gender).Required().Errors(),
			fields.Str("Password", s.Password).Required().Errors(),
			fields.Struct("MainDept", s.MainDept),
			fields.Slice[struct{}]("UserDepts", make([]struct{}, len(s.UserDepts))).Required().Errors(),
			fields.SliceStruct("UserDepts", s.UserDepts),
		)
	case "Update":
		return validatex.Check(
			fields.Int("ID", s.ID).Required().Gt(0).Errors(),
			fields.Str("Username", s.Username).Required().Errors(),
			fields.Str("FullName", s.FullName).Required().Errors(),
			fields.Str("Nickname", s.Nickname).Required().Errors(),
			fields.Int("Gender", s.Gender).Required().Errors(),
			fields.Struct("MainDept", s.MainDept),
			fields.Slice[struct{}]("UserDepts", make([]struct{}, len(s.UserDepts))).Required().Errors(),
			fields.SliceStruct("UserDepts", s.UserDepts),
		)
	case "Detail", "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	case "BatchDelete":
		return fields.Slice("IDs", s.IDs).Required().
			Each("gt", "", func(value int64) bool {
				return value > 0
			}, int64(0)).
			Errors()
	case "Import":
		return validatex.Check(
			fields.Slice[struct{}]("Data", make([]struct{}, len(s.Data))).Required().Min(1).Errors(),
			fields.SliceStruct("Data", s.Data),
		)
	case "Password":
		return validatex.Check(
			fields.Int("ID", s.ID).Required().Gt(0).Errors(),
			fields.Str("Password", s.Password).Required().Errors(),
		)
	default:
		return nil
	}
}

// Validate 部门验证
func (s Dept) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Int("DepartmentId", s.DepartmentId).Required().Gt(0).Errors()
}

// Validate 用户导入子项验证
func (s UserImportItem) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Str("Username", s.Username).Required().Min(3).Max(20).Regexp(usernamePattern).Errors(),
		fields.Str("Password", s.Password).Required().Errors(),
		fields.Str("FullName", s.FullName).Required().Errors(),
		fields.Str("Nickname", s.Nickname).Required().Errors(),
		fields.Str("Email", s.Email).Required().Email().Errors(),
		fields.Int("Gender", s.Gender).Required().Errors(),
	)
}

// Messages 用户请求错误消息
func (s User) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户请求字段翻译
func (s User) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户创建错误消息
func (s UserCreate) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户创建字段翻译
func (s UserCreate) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户更新错误消息
func (s UserUpdate) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户更新字段翻译
func (s UserUpdate) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户密码错误消息
func (s UserPassword) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户密码字段翻译
func (s UserPassword) Translates() validatex.Translates {
	return userTranslates()
}

// Messages 用户导入子项错误消息
func (s UserImportItem) Messages() validatex.Messages {
	return userImportMessages()
}

// Translates 用户导入子项字段翻译
func (s UserImportItem) Translates() validatex.Translates {
	return validatex.Translates{
		"Username": "用户名",
		"Password": "密码",
		"FullName": "姓名",
		"Nickname": "昵称",
		"Email":    "邮箱",
		"Gender":   "性别",
		"Age":      "年龄",
		"Status":   "状态",
	}
}

// Messages 用户导入错误消息
func (s UserImport) Messages() validatex.Messages {
	return userImportMessages()
}

// Translates 用户导入字段翻译
func (s UserImport) Translates() validatex.Translates {
	return validatex.Translates{
		"Data":            "导入数据",
		"Data.*.Username": "用户名",
		"Data.*.Password": "密码",
		"Data.*.FullName": "姓名",
		"Data.*.Nickname": "昵称",
		"Data.*.Email":    "邮箱",
		"Data.*.Gender":   "性别",
		"Data.*.Age":      "年龄",
		"Data.*.Status":   "状态",
	}
}

// Messages 用户批量删除错误消息
func (s UserBatchDelete) Messages() validatex.Messages {
	return userMessages()
}

// Translates 用户批量删除字段翻译
func (s UserBatchDelete) Translates() validatex.Translates {
	return userTranslates()
}

// userMessages 返回用户校验消息
func userMessages() validatex.Messages {
	return validatex.Messages{
		"required":            "字段 {field} 必填",
		"gt":                  "字段 {field} 需大于 0",
		"IDs.*.gt":            "字段 {field} 需大于 0",
		"min":                 "{field} 长度不能少于 {min} 个字符",
		"max":                 "{field} 长度不能超过 {max} 个字符",
		"regex":               "{field} 格式错误",
		"email":               "{field} 邮箱格式错误",
		"Data.min":            "{field} 数量不能少于 {min}",
		"Data.*.Username.min": "{field} 长度不能少于 {min} 个字符",
		"Data.*.Username.max": "{field} 长度不能超过 {max} 个字符",
	}
}

// userImportMessages 返回用户导入校验消息
func userImportMessages() validatex.Messages {
	return validatex.Messages{
		"required":            "字段 {field} 必填",
		"min":                 "{field} 长度不能少于 {min} 个字符",
		"max":                 "{field} 长度不能超过 {max} 个字符",
		"regex":               "{field} 格式错误",
		"email":               "{field} 邮箱格式错误",
		"Data.min":            "{field} 数量不能少于 {min}",
		"Data.*.Username.min": "{field} 长度不能少于 {min} 个字符",
		"Data.*.Username.max": "{field} 长度不能超过 {max} 个字符",
	}
}

// userTranslates 返回用户字段翻译
func userTranslates() validatex.Translates {
	return validatex.Translates{
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
		"ID":                        "ID",
		"Username":                  "用户名",
		"FullName":                  "姓名",
		"Nickname":                  "昵称",
		"Gender":                    "性别",
		"Password":                  "密码",
		"Age":                       "年龄",
		"MainDept.DepartmentId":     "主部门",
		"UserDepts":                 "用户部门",
		"UserDepts.*.DepartmentId":  "用户部门ID",
		"IDs":                       "ID列表",
		"IDs.*":                     "ID",
		"Data":                      "导入数据",
		"Data.*.Username":           "用户名",
		"Data.*.Password":           "密码",
		"Data.*.FullName":           "姓名",
		"Data.*.Nickname":           "昵称",
		"Data.*.Email":              "邮箱",
		"Data.*.Gender":             "性别",
		"Data.*.Age":                "年龄",
		"Data.*.Status":             "状态",
	}
}
