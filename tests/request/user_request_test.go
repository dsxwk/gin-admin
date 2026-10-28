package request

import (
	"gin/app/facade"
	"gin/app/request"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUserListPageValidation 验证用户列表页码条件必填
func TestUserListPageValidation(t *testing.T) {
	t.Run("page required when paginated", func(t *testing.T) {
		ctx := newRequestTestContext("GET", "/?pageSize=10", "", "")
		var req request.User

		err := facade.Request().BindValidate(ctx, &req, "List")
		errors := requireValidationErrors(t, err)
		require.Equal(t, []string{"page", "page"}, validationFieldPaths(errors))
		require.Equal(t, []string{"required", "gt"}, validationErrorRules(errors))
	})

	t.Run("page size required when paginated", func(t *testing.T) {
		ctx := newRequestTestContext("GET", "/?page=1", "", "")
		var req request.User

		err := facade.Request().BindValidate(ctx, &req, "List")
		errors := requireValidationErrors(t, err)
		require.Equal(t, []string{"pageSize", "pageSize"}, validationFieldPaths(errors))
		require.Equal(t, []string{"required", "gt"}, validationErrorRules(errors))
	})

	t.Run("page settings optional when not paginated", func(t *testing.T) {
		ctx := newRequestTestContext("GET", "/?notPage=true", "", "")
		var req request.User

		require.NoError(t, facade.Request().BindValidate(ctx, &req, "List"))
	})
}

// TestUserBatchDeleteIDValidation 验证用户批量删除ID逐项校验
func TestUserBatchDeleteIDValidation(t *testing.T) {
	ctx := newRequestTestContext("POST", "/", "", "")

	tests := []struct {
		name string
		data any
	}{
		{name: "user", data: &request.User{IDs: []int64{1, 0, -1}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := facade.Request().Validate(ctx, test.data, "BatchDelete")
			errors := requireValidationErrors(t, err)
			require.Equal(t, []string{"ids[1]", "ids[2]"}, validationFieldPaths(errors))
			require.Equal(t, "字段 ID 需大于 0", errors[0].Msg)
			require.Equal(t, "字段 ID 需大于 0", errors[1].Msg)
		})
	}
}

// TestUserNestedValidationFields 验证用户嵌套字段路径
func TestUserNestedValidationFields(t *testing.T) {
	ctx := newRequestTestContext("POST", "/", "", "")
	req := request.User{
		Username:  "tester",
		FullName:  "测试用户",
		Nickname:  "tester",
		Gender:    1,
		Password:  "123456",
		UserDepts: []request.Dept{{}},
	}

	err := facade.Request().Validate(ctx, &req, "Create")
	errors := requireValidationErrors(t, err)
	paths := validationFieldPaths(errors)
	require.Contains(t, paths, "mainDept.departmentId")
	require.Contains(t, paths, "userDepts[0].departmentId")
}

// TestUserImportArrayValidationFields 验证用户导入数组字段路径
func TestUserImportArrayValidationFields(t *testing.T) {
	ctx := newRequestTestContext("POST", "/", "", "")
	req := request.User{
		Data: []request.UserImportItem{{Username: "a"}},
	}

	err := facade.Request().Validate(ctx, &req, "Import")
	errors := requireValidationErrors(t, err)
	paths := validationFieldPaths(errors)
	require.Contains(t, paths, "data[0].username")
	require.Contains(t, paths, "data[0].password")
}

// TestMenuNestedValidationFields 验证菜单元数据字段路径
func TestMenuNestedValidationFields(t *testing.T) {
	ctx := newRequestTestContext("POST", "/", "", "")
	req := request.Menu{
		Type:   1,
		Name:   "dashboard",
		Status: 1,
	}

	err := facade.Request().Validate(ctx, &req, "Create")
	errors := requireValidationErrors(t, err)
	paths := validationFieldPaths(errors)
	require.Contains(t, paths, "meta.title")
	require.Contains(t, paths, "meta.component")
}
