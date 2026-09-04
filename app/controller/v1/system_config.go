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

type SystemConfigController struct {
	base.BaseController
	service service.SystemConfigService
}

// List 列表
// @Tags 系统配置管理
// @Summary 列表
// @Description 系统配置列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.SystemConfig}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/system-config [get]
func (s *SystemConfigController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.SystemConfig
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

// UpdateConfig 保存配置
// @Tags 系统配置管理
// @Summary 保存配置
// @Description 保存配置
// @Param token header string true "认证Token"
// @Param data body request.SystemConfigUpdates true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
func (s *SystemConfigController) UpdateConfig(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.SystemConfigUpdates
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

	err = req.Validate(ctx)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	err = s.service.UpdateConfig(ctx, data)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(data))
}

// Detail 详情
// @Tags 系统配置管理
// @Summary 详情
// @Description 系统配置详情
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.SystemConfig} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/system-config/{id} [get]
func (s *SystemConfigController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.SystemConfig
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

// Create 创建
// @Tags 系统配置管理
// @Summary 创建
// @Description 系统配置创建
// @Param token header string true "认证Token"
// @Param data body request.SystemConfigCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.SystemConfig} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/system-config [post]
func (s *SystemConfigController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.SystemConfig
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

// Update 更新
// @Tags 系统配置管理
// @Summary 更新
// @Description 系统配置更新
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Param data body request.SystemConfigUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/system-config/{id} [put]
func (s *SystemConfigController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.SystemConfig
	)

	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		s.Response.Error(c, err)
		return
	}
	err = mapstructure.Decode(data, &req)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
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

// Delete 删除
// @Tags 系统配置管理
// @Summary 删除
// @Description 系统配置删除
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/system-config/{id} [delete]
func (s *SystemConfigController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.SystemConfig
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
