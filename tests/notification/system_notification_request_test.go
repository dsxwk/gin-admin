package notification

import (
	"gin/app/request"
	"testing"

	"github.com/dsxwk/validatex"
	"github.com/stretchr/testify/require"
)

// TestSystemNotificationPushValidate 推送通知请求验证测试
func TestSystemNotificationPushValidate(t *testing.T) {
	req := request.SystemNotificationPush{
		Title:   " 系统通知 ",
		Content: " 通知内容 ",
		Type:    1,
		IsToAll: 2,
		UserIds: []int64{1, 2},
	}
	require.Empty(t, validatex.ValidateSceneRequest(&req, "Push"))
	require.Equal(t, "系统通知", req.Title)
	require.Equal(t, "通知内容", req.Content)

	allReq := request.SystemNotificationPush{
		Title:   "全站通知",
		Content: "通知内容",
		IsToAll: 1,
	}
	require.Empty(t, validatex.ValidateSceneRequest(&allReq, "Push"))

	invalidUserReq := request.SystemNotificationPush{
		Title:   "系统通知",
		Content: "通知内容",
		Type:    1,
		IsToAll: 2,
		UserIds: []int64{0},
	}
	require.NotEmpty(t, validatex.ValidateSceneRequest(&invalidUserReq, "Push"))

	invalidTypeReq := request.SystemNotificationPush{
		Title:   "系统通知",
		Content: "通知内容",
		Type:    3,
		IsToAll: 2,
		UserIds: []int64{1},
	}
	require.NotEmpty(t, validatex.ValidateSceneRequest(&invalidTypeReq, "Push"))
}

// TestSystemNotificationListValidate 通知列表请求验证测试
func TestSystemNotificationListValidate(t *testing.T) {
	req := request.SystemNotificationList{}
	require.Empty(t, validatex.ValidateSceneRequest(&req, "List"))
	require.Equal(t, 1, req.Page)
	require.Equal(t, 20, req.PageSize)
}

// TestSystemNotificationReadValidate 通知已读请求验证测试
func TestSystemNotificationReadValidate(t *testing.T) {
	req := request.SystemNotificationRead{IDs: []int64{1, 2}}
	require.Empty(t, validatex.ValidateSceneRequest(&req, "Read"))

	invalidReq := request.SystemNotificationRead{IDs: []int64{1, 0}}
	require.NotEmpty(t, validatex.ValidateSceneRequest(&invalidReq, "Read"))
}

// TestSystemNotificationRevokeValidate 通知撤回请求验证测试
func TestSystemNotificationRevokeValidate(t *testing.T) {
	req := request.SystemNotificationRevoke{ID: 1}
	require.Empty(t, validatex.ValidateSceneRequest(&req, "Revoke"))

	invalidReq := request.SystemNotificationRevoke{}
	require.NotEmpty(t, validatex.ValidateSceneRequest(&invalidReq, "Revoke"))
}
