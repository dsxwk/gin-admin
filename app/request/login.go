package request

import (
	"fmt"
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// UserLogin 用户登录
type UserLogin struct {
	CaptchaID string `json:"captchaId" label:"验证码ID"`
	Code      string `json:"code" label:"验证码"`
	Username  string `json:"username" example:"admin" label:"用户名"`
	Password  string `json:"password" example:"123456" label:"密码"`
}

// RefreshToken 刷新token
type RefreshToken struct {
	Token string `json:"token" label:"刷新令牌"`
}

type CheckCaptcha struct {
	CaptchaID string `json:"captchaId" label:"验证码ID"`
	Code      string `json:"code" label:"验证码"`
}

// Register 用户注册
type Register struct{}

// Login Validator
type Login struct {
	base.Request
	UserLogin
	RefreshToken RefreshToken
}

// ValidateWithScene 请求验证
func (s Login) ValidateWithScene(scene string, translator *validatex.Translator) validatex.Errors {
	switch scene {
	case "Login":
		return validatex.Check(
			validatex.Str(translator, "UserLogin.Username", s.Username).Required().Errors(),
			validatex.Str(translator, "UserLogin.Password", s.Password).Required().Errors(),
		)
	case "RefreshToken":
		return validatex.Str(translator, "RefreshToken.Token", s.RefreshToken.Token).Required().Errors()
	case "CheckCaptcha":
		return validatex.Check(
			validatex.Str(translator, "UserLogin.CaptchaID", s.CaptchaID).Required().Errors(),
			validatex.Str(translator, "UserLogin.Code", s.Code).Required().Errors(),
		)
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Login) Messages() validatex.Messages {
	return validatex.Messages{
		"required": fmt.Sprintf("%s {field} %s",
			s.Trans("validator.common.field", nil),
			s.Trans("validator.common.required", nil),
		),
	}
}

// Translates 字段翻译
func (s Login) Translates() validatex.Translates {
	return validatex.Translates{
		"UserLogin.Username":  s.Trans("validator.login.username", nil),
		"UserLogin.Password":  s.Trans("validator.login.password", nil),
		"RefreshToken.Token":  s.Trans("validator.login.refreshToken", nil),
		"UserLogin.CaptchaID": s.Trans("login.captchaId", nil),
		"UserLogin.Code":      s.Trans("login.code", nil),
	}
}
