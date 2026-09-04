package service

import (
	"context"
	"errors"
	"gin/app/event"
	"gin/app/facade"
	"gin/app/middleware"
	"gin/app/model"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg"

	"gorm.io/gorm"
)

type LoginService struct {
	base.BaseService
}

// Login 登录
func (s *LoginService) Login(ctx context.Context, username, password string) (err error, m model.User, accessToken, refreshToken string, tokenExpire, refreshTokenExpire int64) {
	var (
		db   = s.DB(ctx, &m)
		conf = facade.Config()
	)

	if err = db.Model(&m).
		Preload("UserRoles").
		Where("username = ?", username).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.NotFound().WithMsg("login.accountErr"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
		}
	}

	check := pkg.BcryptCheck(password, m.Password)
	if !check {
		return errcode.ArgsError().WithMsg("login.pwdErr"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	if m.Status == 2 {
		return errcode.ArgsError().WithMsg("login.accountDisabled"), m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	accessToken, refreshToken, tokenExpire, refreshTokenExpire, err = middleware.Jwt{}.WithRefresh(m.ID, conf.Jwt.Exp, conf.Jwt.RefreshExp)
	if err != nil {
		return err, m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
	}

	// 发布事件
	facade.Event().Publish[event.UserLoginEvent](ctx, event.UserLoginEvent{
		UserId:   m.ID,
		Username: m.Username,
	})

	return nil, m, accessToken, refreshToken, tokenExpire, refreshTokenExpire
}

// RefreshToken 刷新token
func (s *LoginService) RefreshToken(ctx context.Context, token string) (accessToken, refreshToken string, tExp, rExp int64, err error) {
	var (
		conf = facade.Config()
	)

	j := middleware.Jwt{}
	claims, err := j.Decode(token)
	if err != nil || claims["typ"] != "refresh" {
		return accessToken, refreshToken, tExp, rExp, errcode.ArgsError().WithMsg("login.invalidToken")
	}

	uid := int64(claims["id"].(float64))

	return j.WithRefresh(uid, conf.Jwt.Exp, conf.Jwt.RefreshExp)
}
