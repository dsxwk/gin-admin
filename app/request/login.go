package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// UserLogin 用户登录
type UserLogin struct {
	CaptchaID string `json:"captchaId" form:"captchaId" label:"验证码ID"`
	Code      string `json:"code" form:"code" label:"验证码"`
	Username  string `json:"username" form:"username" example:"admin" label:"用户名"`
	Password  string `json:"password" form:"password" example:"123456" label:"密码"`
}

// RefreshToken 刷新token
type RefreshToken struct {
	Token string `json:"token" form:"token" label:"刷新令牌"`
}

// CheckCaptcha 验证码校验
type CheckCaptcha struct {
	CaptchaID string `json:"captchaId" form:"captchaId" label:"验证码ID"`
	Code      string `json:"code" form:"code" label:"验证码"`
}

// Register 用户注册
type Register struct{}

// Login 登录验证
type Login struct {
	base.Request
	UserLogin
	RefreshToken RefreshToken
}

// ValidateWithScene 请求验证
func (s Login) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Login":
		return validatex.Check(
			fields.Str("username", s.Username).Required().Errors(),
			fields.Str("password", s.Password).Required().Errors(),
		)
	case "RefreshToken":
		return fields.Str("token", s.RefreshToken.Token).Required().Errors()
	case "CheckCaptcha":
		return validatex.Check(
			fields.Str("captchaId", s.CaptchaID).Required().Errors(),
			fields.Str("code", s.Code).Required().Errors(),
		)
	default:
		return nil
	}
}

// Messages 错误消息
func (s Login) Messages() validatex.Messages {
	return validatex.Messages{
		"required": s.Trans("validator.common.requiredField", nil),
	}
}

// Translates 字段翻译
func (s Login) Translates() validatex.Translates {
	return validatex.Translates{
		"username":  s.Trans("validator.login.username", nil),
		"password":  s.Trans("validator.login.password", nil),
		"token":     s.Trans("validator.login.refreshToken", nil),
		"captchaId": s.Trans("login.captchaId", nil),
		"code":      s.Trans("login.code", nil),
	}
}
