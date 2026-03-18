package service

import (
	"errors"
	"gin/app/event"
	"gin/app/middleware"
	"gin/app/model"
	"gin/common/base"
	"gin/pkg"
	"gin/pkg/container"
	"gin/pkg/eventbus"
	"gorm.io/gorm"
)

type LoginService struct {
	base.BaseService
}

// Login 登录
func (s *LoginService) Login(username, password string) (err error, m model.User, accessToken, refreshToken string, tokenExpire, refreshTokenExpire int64) {
	containers := container.Get(s.Ctx)
	if err = containers.DB.
		Where("username = ?", username).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("login.accountErr"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
		}
	}

	check := pkg.BcryptCheck(password, m.Password)
	if !check {
		return errors.New("login.pwdErr"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	if m.Status == 2 {
		return errors.New("login.accountDisabled"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	jwt := middleware.Jwt{}
	accessToken, refreshToken, tokenExpire, refreshTokenExpire, err = jwt.WithRefresh(m.ID, containers.Config.Jwt.Exp, containers.Config.Jwt.RefreshExp)
	if err != nil {
		return errors.New(err.Error()), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	// 发布事件
	eventbus.Publish(s.Ctx, event.UserLoginEvent{
		UserId:   m.ID,
		Username: m.Username,
	})

	return nil, m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
}

// RefreshToken 刷新token
func (s *LoginService) RefreshToken(token string) (accessToken, refreshToken string, tExp, rExp int64, err error) {
	containers := container.Get(s.Ctx)
	j := middleware.Jwt{}
	claims, err := j.Decode(token)
	if err != nil || claims["typ"] != "refresh" {
		return accessToken, refreshToken, tExp, rExp, errors.New("login.invalidToken")
	}

	uid := int64(claims["id"].(float64))

	return j.WithRefresh(uid, containers.Config.Jwt.Exp, containers.Config.Jwt.RefreshExp)
}
