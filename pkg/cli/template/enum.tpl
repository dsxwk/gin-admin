package {{.Package}}

import (
	enumx "gin/pkg/enum"
)

// {{.Name}}Enum {{.Description}}
type {{.Name}}Enum struct {}

const (
	UserGenderSecret = 0 // 保密
	UserGenderMale   = 1 // 男
	UserGenderFemale = 2 // 女
)

const (
	UserStatusEnabled  = "enable"  // 启用
	UserStatusDisabled = "disable" // 停用
)

// Gender 性别
func (s *{{.Name}}Enum) Gender() *enumx.Enum[int] {
	return enumx.NewEnum(
		enumx.Item[int]{Value: UserGenderSecret, Desc: "保密"},
		enumx.Item[int]{Value: UserGenderMale, Desc: "男"},
		enumx.Item[int]{Value: UserGenderFemale, Desc: "女"},
	)
}

// Status 状态
func (s *{{.Name}}Enum) Status() *enumx.Enum[string] {
	return enumx.NewEnum(
		enumx.Item[string]{Value: UserStatusEnabled, Desc: "启用"},
		enumx.Item[string]{Value: UserStatusDisabled, Desc: "停用"},
	)
}
