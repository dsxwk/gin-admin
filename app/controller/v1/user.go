package v1

import (
	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-viper/mapstructure/v2"
)

type UserController struct {
	base.BaseController
	service service.UserService
}

// List 列表
// @Tags 用户管理
// @Summary 列表
// @Description 用户列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.User}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user [get]
func (s *UserController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.User
	)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "List")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	res, err := s.service.List(ctx, req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// Create 创建
// @Tags 用户管理
// @Summary 创建
// @Description 用户创建
// @Param token header string true "认证Token"
// @Param data body request.UserCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user [post]
func (s *UserController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.User
	)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Create")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	user, err := s.service.Create(ctx, req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(user))
}

// Update 更新
// @Tags 用户管理
// @Summary 更新
// @Description 用户更新
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Param data body request.UserUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [put]
func (s *UserController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.User
	)

	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	err = mapstructure.Decode(data, &req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	req.ID = facade.Request().Path[int64](c, "id", 0)
	err = facade.Request().Validate(c, &req, "Update")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	err = s.service.Update(ctx, req.ID, data)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(data))
}

// Detail 详情
// @Tags 用户管理
// @Summary 详情
// @Description 用户详情
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [get]
func (s *UserController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.User
	)

	req.ID = facade.Request().Path[int64](c, "id", 0)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Detail")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	m, err := s.service.Detail(ctx, req.ID)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Delete 删除
// @Tags 用户管理
// @Summary 删除
// @Description 用户删除
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [delete]
func (s *UserController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.User
	)

	req.ID = facade.Request().Path[int64](c, "id", 0)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Delete")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	err = s.service.Delete(ctx, req.ID)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success())
}

// Import 批量导入用户
// @Tags 用户管理
// @Summary 批量导入
// @Description 批量导入用户
// @Param token header string true "认证Token"
// @Param data body request.UserImport true "导入参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/import [post]
func (s *UserController) Import(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.UserImport
	)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Import")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	result, err := s.service.Import(ctx, s.GetUserId(c), req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(result))
}

// BatchDelete 批量删除
// @Tags 用户管理
// @Summary 批量删除
// @Description 用户批量删除
// @Param token header string true "认证Token"
// @Param data body request.UserBatchDelete true "批量删除参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/batch-delete [post]
func (s *UserController) BatchDelete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.User
	)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "BatchDelete")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	err = s.service.BatchDelete(ctx, req.IDs)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success())
}

// Password 更新密码
// @Tags 用户管理
// @Summary 更新密码
// @Description 更新密码
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Param data body request.UserPassword true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id}/password [put]
func (s *UserController) Password(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		req      request.User
		userPass request.UserPassword
	)

	req.ID = facade.Request().Path[int64](c, "id", 0)

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Password")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	userPass.Password = req.Password

	err = s.service.Password(ctx, req.ID, userPass)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success())
}
