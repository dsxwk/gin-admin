package request

import (
	"gin/app/facade"
	"gin/app/request"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginRequestValidation(t *testing.T) {
	t.Run("login keeps captcha optional", func(t *testing.T) {
		ctx := newRequestTestContext(
			"POST",
			"/",
			`{"username":"admin","password":"123456"}`,
			"application/json",
		)
		var req request.Login

		require.NoError(t, facade.Request().BindValidate(ctx, &req, "Login"))
	})

	t.Run("login requires credentials", func(t *testing.T) {
		ctx := newRequestTestContext("POST", "/", `{}`, "application/json")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "Login")
		require.Error(t, err)
		require.ErrorContains(t, err, "UserLogin.Username")
		require.ErrorContains(t, err, "UserLogin.Password")
	})

	t.Run("refresh token requires token", func(t *testing.T) {
		ctx := newRequestTestContext("POST", "/", "", "")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "RefreshToken")
		require.Error(t, err)
		require.ErrorContains(t, err, "RefreshToken.Token")

		req.RefreshToken.Token = "refresh-token"
		require.NoError(t, facade.Request().Validate(ctx, &req, "RefreshToken"))
	})

	t.Run("captcha check requires captcha", func(t *testing.T) {
		ctx := newRequestTestContext("POST", "/", `{}`, "application/json")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "CheckCaptcha")
		require.Error(t, err)
		require.ErrorContains(t, err, "UserLogin.CaptchaID")
		require.ErrorContains(t, err, "UserLogin.Code")

		req.CaptchaID = "captcha-id"
		req.Code = "code"
		require.NoError(t, facade.Request().Validate(ctx, &req, "CheckCaptcha"))
	})
}
