package enum

import (
	"gin/common/base"
)

const (
	ImportRecordsTypeUser = 1 // 用户导入
)

// ImportRecordsEnum 导入记录枚举
type ImportRecordsEnum struct{}

// Type 导入类型
func (s *ImportRecordsEnum) Type() *base.Enum[int] {
	return base.NewEnum(
		base.Item[int]{Value: ImportRecordsTypeUser, Desc: "用户导入"},
	)
}
