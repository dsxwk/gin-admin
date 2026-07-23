package v1

import (
	"gin/app/enum"
	"gin/app/facade"
	"gin/app/job"
	"gin/app/model"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg/serviceprovider/lang"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
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
	Token Token      `json:"token"`
	User  model.User `json:"user"`
}

type CaptchaResponse struct {
	CaptchaId    string `json:"captchaId"`
	CaptchaImage string `json:"captchaImage"`
}

// Login 登录
// @Tags 登录相关
// @Summary 账号密码登录
// @Description 用户账号密码登录
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

	req.WithContext(ctx)
	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Login")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	if req.Code != "" {
		ok := s.verify(req.CaptchaID, req.Code)
		if !ok {
			s.Response.Error(c, errcode.ArgsError().WithMsg("验证码错误"))
			return
		}
	}

	err, userModel, accessToken, refreshToken, tokenExpire, refreshTokenExpire := s.service.Login(req.Username, req.Password)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(
		c, errcode.Success().WithMsg(
			facade.Lang().Trans(ctx, "login.success", map[string]interface{}{
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

	req.WithContext(ctx)
	s.service.WithContext(ctx)

	token := c.Request.Header.Get("token")
	req.RefreshToken.Token = token
	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "RefreshToken")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	accessToken, refreshToken, tokenExpire, refreshTokenExpire, err := s.service.RefreshToken(token)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		TokenExpire:        tokenExpire,
		RefreshTokenExpire: refreshTokenExpire,
	}))
}

// Test 测试
// @Tags 登录相关
// @Summary 测试
// @Description 测试
// @Accept json
// @Produce json
// @Success 200 {object} errcode.SuccessResponse{data=map[string]any{}} "成功"
// @Router /api/v1/test [post]
func (s *LoginController) Test(c *gin.Context) {
	var (
		userEnum enum.UserEnum
		ctx      = c.Request.Context()
	)

	status := userEnum.Status().Get()
	desc1 := userEnum.Status().Desc(enum.UserStatusEnabled)
	value1 := userEnum.Status().Value("启用")
	_map := userEnum.Status().Map()
	containsValue := userEnum.Status().ContainsValue(enum.UserStatusEnabled)
	containsDesc := userEnum.Status().ContainsDesc("启用")
	length := userEnum.Status().Len()

	gender := userEnum.Gender().Get()
	desc2 := userEnum.Gender().Desc(enum.UserGenderMale)
	value2 := userEnum.Gender().Value("男")
	_map2 := userEnum.Gender().Map()
	containsValue2 := userEnum.Gender().ContainsValue(enum.UserGenderMale)
	containsDesc2 := userEnum.Gender().ContainsDesc("男")
	length2 := userEnum.Gender().Len()

	_ = facade.Queue().Producer("kafka_demo").Publish(ctx, []byte(`{"name":"kafka_test111"}`))
	_ = facade.Queue().Producer("kafka_delay_demo").Publish(ctx, []byte(`{"name":"kafka_test222"}`))
	_ = facade.Queue().Producer("rabbitmq_demo").Publish(ctx, []byte(`{"name":"test111"}`))
	_ = facade.Queue().Producer("rabbitmq_delay_demo").Publish(ctx, []byte(`{"name":"test222"}`))
	_ = facade.Job().Dispatch(ctx, "send_email", job.SendEmail{
		To:      "a@b.com",
		Subject: "你好",
		Content: "这是一封测试邮件",
	})
	_ = facade.Job().Dispatch(ctx, "export_report", job.ExportReport{ReportType: "daily", UserID: 1})
	_ = facade.Job().Dispatch(ctx, "sync_user", job.SyncUser{UserID: 1, Action: "update"})

	s.Response.Success(c, errcode.Success().WithData(map[string]any{
		"status":         status,
		"desc1":          desc1,
		"value1":         value1,
		"map":            _map,
		"containsValue":  containsValue,
		"containsDesc":   containsDesc,
		"length":         length,
		"gender":         gender,
		"desc2":          desc2,
		"value2":         value2,
		"map2":           _map2,
		"containsValue2": containsValue2,
		"containsDesc2":  containsDesc2,
		"length2":        length2,
	}))
}

// generateCharset 生成字符串
func (s *LoginController) generateCharset() string {
	charset := ""

	// 数字0-9
	for i := '0'; i <= '9'; i++ {
		charset += string(i)
	}

	// 小写字母a-z
	for i := 'a'; i <= 'z'; i++ {
		charset += string(i)
	}

	// 大写字母A-Z
	for i := 'A'; i <= 'Z'; i++ {
		charset += string(i)
	}

	return charset
}

// GetCaptcha 获取验证码
// @Tags 登录相关
// @Summary 获取验证码
// @Description 获取验证码
// @Accept json
// @Produce json
// @Success 200 {object} errcode.SuccessResponse{data=CaptchaResponse} "成功返回" Example({"code":0,"msg":"Success","data":[]})
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误" Example({"code":500,"msg":"系统错误","data":[]})
// @Router /api/v1/captcha [get]
func (s *LoginController) GetCaptcha(c *gin.Context) {
	// 配置验证码
	driver := base64Captcha.NewDriverString(
		60,                  // 高度
		240,                 // 宽度
		0,                   // 噪点数量
		0,                   // base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSlimeLine, // 显示线条选项
		4,                   // 验证码长度
		s.generateCharset(), // 验证码字符集
		&color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 背景颜色(白色)
		base64Captcha.DefaultEmbeddedFonts,          // 字体存储
		nil,                                         // []string{"wqy-microfiche.ttf"},              // 字体名称
	)

	// 生成验证码
	id, b64s, _, err := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore).Generate()
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	// 返回验证码ID和图片的Base64数据
	s.Response.Success(c, errcode.Success().WithData(CaptchaResponse{
		CaptchaId:    id,
		CaptchaImage: b64s,
	}))
}

// CheckCaptcha 校验验证码
// @Tags 登录相关
// @Summary 校验验证码
// @Description 校验验证码
// @Param data body request.CheckCaptcha true "验证码校验参数"
// @Success 200 {object} errcode.SuccessResponse{data=Token} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Router /api/v1/captcha [post]
func (s *LoginController) CheckCaptcha(c *gin.Context) {
	var (
		req request.Login
	)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "CheckCaptcha")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(s.verify(req.CaptchaID, req.Code)))
}

// verify 验证
func (s *LoginController) verify(captchaId, value string) bool {
	return base64Captcha.DefaultMemStore.Verify(captchaId, value, true)
}
