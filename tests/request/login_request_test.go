package request

import (
	"encoding/json"
	"gin/app/facade"
	"gin/app/request"
	"net/http"
	"testing"

	"github.com/dsxwk/validatex"
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
		errors := requireValidationErrors(t, err)
		require.Equal(t, []string{"username", "password"}, validationFieldPaths(errors))
	})

	t.Run("refresh token requires token", func(t *testing.T) {
		ctx := newRequestTestContext("POST", "/", "", "")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "RefreshToken")
		require.Error(t, err)
		errors := requireValidationErrors(t, err)
		require.Equal(t, []string{"token"}, validationFieldPaths(errors))

		req.RefreshToken.Token = "refresh-token"
		require.NoError(t, facade.Request().Validate(ctx, &req, "RefreshToken"))
	})

	t.Run("captcha check requires captcha", func(t *testing.T) {
		ctx := newRequestTestContext("POST", "/", `{}`, "application/json")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "CheckCaptcha")
		require.Error(t, err)
		errors := requireValidationErrors(t, err)
		require.Equal(t, []string{"captchaId", "code"}, validationFieldPaths(errors))

		req.CaptchaID = "captcha-id"
		req.Code = "code"
		require.NoError(t, facade.Request().Validate(ctx, &req, "CheckCaptcha"))
	})

	t.Run("validation error returns structured fields", func(t *testing.T) {
		ctx, recorder := newRequestTestContextWithRecorder("POST", "/", `{}`, "application/json")
		var req request.Login

		err := facade.Request().BindValidate(ctx, &req, "Login")
		require.Error(t, err)

		facade.Response().Error(ctx, err)
		require.Equal(t, http.StatusBadRequest, ctx.Writer.Status())

		var response struct {
			Code int64                         `json:"code"`
			Msg  string                        `json:"msg"`
			Data validatex.ValidationErrorData `json:"data"`
		}
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		require.Equal(t, int64(400), response.Code)
		require.Equal(t, "用户名不能为空", response.Msg)
		require.Len(t, response.Data.Errors, 2)
		require.Equal(t, "username", response.Data.Errors[0].Path)
		require.Equal(t, "required", response.Data.Errors[0].Rule)
		require.Equal(t, "用户名不能为空", response.Data.Errors[0].Msg)
		require.Equal(t, "password", response.Data.Errors[1].Path)
		require.Equal(t, "required", response.Data.Errors[1].Rule)
		require.Equal(t, "密码不能为空", response.Data.Errors[1].Msg)
	})
}
