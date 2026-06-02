package v1

import (
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg/serviceprovider/lang"
	"github.com/gin-gonic/gin"
)

type MenuController struct {
	base.BaseController
	service service.MenuService
}

// List 列表
// @Tags 菜单管理
// @Summary 列表
// @Description 菜单列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.Menu}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu [get]
func (s *MenuController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "List")
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

// RoleMenu 角色菜单
// @Tags 菜单管理
// @Summary 角色菜单
// @Description 角色菜单
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=[]pkg.TreeNode} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role/{roleId}/menu [get]
func (s *MenuController) RoleMenu(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	s.service.WithContext(ctx)

	req.RoleId = facade.Request[string]().Path(c, "roleId", "0")

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "RoleMenu")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.RoleMenu(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// Detail 详情
// @Tags 菜单管理
// @Summary 详情
// @Description 菜单详情
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id} [get]
func (s *MenuController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Detail")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.Detail(req.Id)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Create 创建菜单
// @Tags 菜单管理
// @Summary 创建菜单
// @Description 用户创建
// @Param token header string true "认证Token"
// @Param data body request.MenuCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=request.MenuCreate} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu [post]
func (s *MenuController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Create")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.Create(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Action 菜单功能
// @Tags 菜单管理
// @Summary 菜单功能
// @Description 菜单功能
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.MenuActions}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/action [get]
func (s *MenuController) Action(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Action")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.Action(req.MenuId)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}
