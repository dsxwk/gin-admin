package enum

import (
	enumx "gin/pkg/enum"
)

const (
	UserGenderSecret = 0 // 保密
	UserGenderMale   = 1 // 男
	UserGenderFemale = 2 // 女
)

const (
	UserStatusEnabled  int64 = 1 // 启用
	UserStatusDisabled int64 = 2 // 停用
)

// UserEnum 用户枚举
type UserEnum struct{}

// Gender 性别
func (s *UserEnum) Gender() *enumx.Enum[int] {
	return enumx.NewEnum(
		enumx.Item[int]{Value: UserGenderSecret, Desc: "保密"},
		enumx.Item[int]{Value: UserGenderMale, Desc: "男"},
		enumx.Item[int]{Value: UserGenderFemale, Desc: "女"},
	)
}

// Status 状态
func (s *UserEnum) Status() *enumx.Enum[int64] {
	return enumx.NewEnum(
		enumx.Item[int64]{Value: UserStatusEnabled, Desc: "启用"},
		enumx.Item[int64]{Value: UserStatusDisabled, Desc: "停用"},
	)
}
