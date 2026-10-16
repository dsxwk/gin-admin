package enum

import enumx "gin/pkg/enum"

// UserDepartmentsEnum 用户部门枚举
type UserDepartmentsEnum struct{}

const (
	DepartmentMainYes int64 = 1
	DepartmentMainNo  int64 = 2
)

// IsMain 是否是主部门
func (s *UserDepartmentsEnum) IsMain() *enumx.Enum[int64] {
	return enumx.NewEnum(
		enumx.Item[int64]{Value: DepartmentMainYes, Desc: "是"},
		enumx.Item[int64]{Value: DepartmentMainNo, Desc: "否"},
	)
}
