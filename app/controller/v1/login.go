package v1

import (
	"gin/app/facade"
	"gin/app/model"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg/lang"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	base.BaseController
	service service.LoginService
}

// Token token信息
type Token struct {
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	TokenExpire        int64  `json:"tokenExpire" example:"7200"`
	RefreshTokenExpire int64  `json:"refreshTokenExpire" example:"172800"`
}

type LoginResponse struct {
	Token Token `json:"token"`
	User  model.User
}

// Login 登录
// @Tags 登录相关
// @Summary 登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param data body request.UserLogin true "登录参数"
// @Success 200 {object} errcode.SuccessResponse{data=LoginResponse} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/login [post]
func (s *LoginController) Login(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Login
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request.BindValidate(c, &req, "Login")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err, userModel, accessToken, refreshToken, tokenExpire, refreshTokenExpire := s.service.Login(req.Username, req.Password)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(lang.T(ctx, err.Error(), nil)))
		return
	}

	s.Success(
		c, errcode.Success().WithMsg(
			facade.Lang.T(ctx, "login.success", map[string]interface{}{
				"name": userModel.Username,
			}),
		).WithData(LoginResponse{
			Token{
				AccessToken:        accessToken,
				RefreshToken:       refreshToken,
				TokenExpire:        tokenExpire,
				RefreshTokenExpire: refreshTokenExpire,
			},
			userModel,
		}),
	)
}

// RefreshToken 刷新token
// @Tags 登录相关
// @Summary 刷新token
// @Description 刷新token
// @Accept json
// @Produce json
// @Param token header string true "刷新Token"
// @Success 200 {object} errcode.SuccessResponse{data=Token} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/refresh-token [post]
func (s *LoginController) RefreshToken(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Login
	)

	s.service.WithContext(ctx)

	token := c.Request.Header.Get("token")
	req.RefreshToken.Token = token
	// 验证
	err := req.Validate(req, "RefreshToken")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	accessToken, refreshToken, tokenExpire, refreshTokenExpire, err := s.service.RefreshToken(token)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(lang.T(ctx, err.Error(), nil)))
		return
	}

	s.Success(c, errcode.Success().WithData(Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		TokenExpire:        tokenExpire,
		RefreshTokenExpire: refreshTokenExpire,
	}))
}

func (s *LoginController) Test(c *gin.Context) {
	s.Success(c, errcode.Success())
}
