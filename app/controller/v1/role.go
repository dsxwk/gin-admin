package v1

import (
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg/serviceprovider/lang"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-viper/mapstructure/v2"
)

type RoleController struct {
	base.BaseController
	service service.RoleService
}

// List 列表
// @Tags 角色管理
// @Summary 列表
// @Description 角色列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.Roles}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role [get]
func (s *RoleController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Roles
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.List(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// Detail 详情
// @Tags 角色管理
// @Summary 详情
// @Description 角色详情
// @Param token header string true "认证Token"
// @Param id path int true "角色ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.Roles} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role/{id} [get]
func (s *RoleController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Roles
	)

	s.service.WithContext(ctx)

	req.ID = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Detail")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.Detail(req.ID)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// Create 创建
// @Tags 角色管理
// @Summary 创建
// @Description 角色创建
// @Param token header string true "认证Token"
// @Param data body request.RoleCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.Roles} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role [post]
func (s *RoleController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Roles
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Create")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	role, err := s.service.Create(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(role))
}

// Update 更新
// @Tags 角色管理
// @Summary 更新
// @Description 角色更新
// @Param token header string true "认证Token"
// @Param id path int true "角色ID"
// @Param data body request.RoleUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role/{id} [put]
func (s *RoleController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]interface{}
		req  request.Roles
	)

	s.service.WithContext(ctx)

	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	err = mapstructure.Decode(data, &req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	req.ID = facade.Request[int64]().Path(c, "id", 0)
	err = req.Validate(req, "Update")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.Update(req.ID, data)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(data))
}

// Delete 删除
// @Tags 角色管理
// @Summary 删除
// @Description 角色删除
// @Param token header string true "认证Token"
// @Param id path int true "角色ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role/{id} [delete]
func (s *RoleController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Roles
	)

	s.service.WithContext(ctx)

	req.ID = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Delete")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.Delete(req.ID)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success())
}
