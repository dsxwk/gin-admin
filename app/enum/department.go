package enum

import enumx "gin/pkg/enum"

const (
	DepartmentStatusEnabled  int64 = 1 // 启用
	DepartmentStatusDisabled int64 = 2 // 停用
)

// DepartmentEnum 部门枚举
type DepartmentEnum struct{}

// Status 状态
func (s *DepartmentEnum) Status() *enumx.Enum[int64] {
	return enumx.NewEnum(
		enumx.Item[int64]{Value: DepartmentStatusEnabled, Desc: "启用"},
		enumx.Item[int64]{Value: DepartmentStatusDisabled, Desc: "停用"},
	)
}
