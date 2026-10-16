package enum

import (
	enumx "gin/pkg/enum"
)

const (
	ImportRecordsTypeUser = 1 // 用户导入
)

// ImportRecordsEnum 导入记录枚举
type ImportRecordsEnum struct{}

// Type 导入类型
func (s *ImportRecordsEnum) Type() *enumx.Enum[int] {
	return enumx.NewEnum(
		enumx.Item[int]{Value: ImportRecordsTypeUser, Desc: "用户导入"},
	)
}
