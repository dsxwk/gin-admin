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

	s.service.WithContext(ctx)

	req.RoleId = facade.Request[string]().Path(c, "id", "0")

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

// Update 更新菜单
// @Tags 菜单管理
// @Summary 创建菜单
// @Description 更新菜单
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param data body request.MenuUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse{data=request.MenuUpdate} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id} [put]
func (s *MenuController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]interface{}
		req  request.Menu
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

	req.Id = facade.Request[int64]().Path(c, "id", 0)
	err = req.Validate(req, "Update")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.Update(req.Id, data)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
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

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Delete")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.Delete(req.Id)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success())
}

// Action 菜单功能
// @Tags 菜单管理
// @Summary 菜单功能
// @Description 菜单功能
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.MenuActions}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/action [get]
func (s *MenuController) Action(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.MenuActions
	)

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "List")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.Action(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(lang.Trans(ctx, err.Error(), nil)))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(res))
}

// ActionDetail 菜单功能详情
// @Tags 菜单管理
// @Summary 菜单功能详情
// @Description 菜单功能详情
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param actionId path int true "菜单功能ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.MenuActions} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/action/{actionId} [get]
func (s *MenuController) ActionDetail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.MenuActions
	)

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "actionId", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Detail")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.ActionDetail(req.Id)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// CreateAction 创建菜单功能
// @Tags 菜单管理
// @Summary 创建菜单功能
// @Description 创建菜单功能
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param data body request.ActionCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=request.ActionCreate} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/action [post]
func (s *MenuController) CreateAction(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.MenuActions
	)

	s.service.WithContext(ctx)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Create")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.CreateAction(req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(req))
}

// UpdateAction 更新菜单功能
// @Tags 菜单管理
// @Summary 更新菜单功能
// @Description 更新菜单功能
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param actionId path int true "菜单功能ID"
// @Param data body request.ActionUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse{data=request.ActionUpdate} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/action/{actionId} [put]
func (s *MenuController) UpdateAction(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]interface{}
		req  request.MenuActions
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

	req.Id = facade.Request[int64]().Path(c, "actionId", 0)
	err = req.Validate(req, "Update")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.UpdateAction(req.Id, data)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(data))
}

// DeleteAction 删除菜单功能
// @Tags 菜单管理
// @Summary 创建菜单功能
// @Description 创建菜单功能
// @Param token header string true "认证Token"
// @Param id path int true "菜单ID"
// @Param actionId path int true "菜单功能ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/menu/{id}/{actionId} [delete]
func (s *MenuController) DeleteAction(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.MenuActions
	)

	s.service.WithContext(ctx)

	req.Id = facade.Request[int64]().Path(c, "actionId", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Delete")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.DeleteAction(req.Id)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success())
}
