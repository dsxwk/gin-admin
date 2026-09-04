package v1

import (
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-viper/mapstructure/v2"
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

// RoleMenu 角色菜单
// @Tags 菜单管理
// @Summary 角色菜单
// @Description 角色菜单
// @Param token header string true "认证Token"
// @Param id path int true "角色ID"
// @Success 200 {object} errcode.SuccessResponse{data=[]pkg.TreeNode} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/role/{id}/menu [get]
func (s *MenuController) RoleMenu(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
	)

	req.RoleId = facade.Request().Path[string](c, "id", "0")

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "RoleMenu")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	res, err := s.service.RoleMenu(ctx, req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// Detail 详情
// @Tags 菜单管理
// @Summary 详情
// @Description 菜单详情
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id} [get]
func (s *MenuController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
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

// Create 创建菜单
// @Tags 菜单管理
// @Summary 创建菜单
// @Description 创建菜单
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

	// 绑定参数并验证
	err := facade.Request().BindValidate(c, &req, "Create")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	m, err := s.service.Create(ctx, req)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Update 更新菜单
// @Tags 菜单管理
// @Summary 创建菜单
// @Description 更新菜单
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param data body request.MenuUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id} [put]
func (s *MenuController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.Menu
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

// Delete 菜单删除
// @Tags 菜单管理
// @Summary 菜单删除
// @Description 菜单删除
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id} [delete]
func (s *MenuController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Menu
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
