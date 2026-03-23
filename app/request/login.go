package request

import (
	"errors"
	"gin/pkg/lang"
	"github.com/gookit/validate"
)

// UserLogin 用户登录
type UserLogin struct {
	Username string `json:"username" validate:"required" example:"admin" label:"用户名"`
	Password string `json:"password" validate:"required" example:"123456" label:"密码"`
}

// RefreshToken 刷新token
type RefreshToken struct {
	Token string `json:"token" validate:"required" label:"刷新令牌"`
}

// Register 用户注册
type Register struct {
}

// Login Validator
type Login struct {
	UserLogin
	RefreshToken RefreshToken
	Context
}

// Validate 请求验证
func (s Login) Validate(data Login, scene string) error {
	v := validate.Struct(data, scene)

	// v.AddMessages(s.Messages())
	// v.AddTranslates(s.Translates())

	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Login) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Login":        []string{"UserLogin.Username", "UserLogin.Password"},
		"RefreshToken": []string{"RefreshToken.Token"},
	})
}

// Messages 验证器错误消息
func (s Login) Messages() map[string]string {
	return validate.MS{
		"required": lang.T(s.Ctx, "validator.common.field", nil) + " {field} " + lang.T(s.Ctx, "validator.common.required", nil),
	}
}

// Translates 字段翻译
func (s Login) Translates() map[string]string {
	return validate.MS{
		"UserLogin.Username": lang.T(s.Ctx, "validator.login.username", nil),
		"UserLogin.Password": lang.T(s.Ctx, "validator.login.password", nil),
		"RefreshToken.Token": lang.T(s.Ctx, "validator.login.refreshToken", nil),
	}
}
