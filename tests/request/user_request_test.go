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
		require.ErrorContains(t, err, "PageListValidate.Page")
	})

	t.Run("page size required when paginated", func(t *testing.T) {
		ctx := newRequestTestContext("GET", "/?page=1", "", "")
		var req request.User

		err := facade.Request().BindValidate(ctx, &req, "List")
		require.ErrorContains(t, err, "PageListValidate.PageSize")
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
			require.ErrorContains(t, err, "IDs[1]")
			require.ErrorContains(t, err, "IDs[2]")
			require.ErrorContains(t, err, "字段 ID 需大于 0")
		})
	}
}
