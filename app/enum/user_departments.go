package enum

import "gin/common/base"

// UserDepartmentsEnum 用户部门枚举
type UserDepartmentsEnum struct{}

const (
	DepartmentMainYes int64 = 1
	DepartmentMainNo  int64 = 2
)

// IsMain 是否是主部门
func (s *UserDepartmentsEnum) IsMain() *base.Enum[int64] {
	return base.NewEnum(
		base.Item[int64]{Value: DepartmentMainYes, Desc: "是"},
		base.Item[int64]{Value: DepartmentMainNo, Desc: "否"},
	)
}
