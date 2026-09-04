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

type DictController struct {
	base.BaseController
	service service.DictService
}

// List 列表
// @Tags 字典管理
// @Summary 列表
// @Description 字典列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.Dict}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dict [get]
func (s *DictController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Dict
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

// Detail 详情
// @Tags 字典管理
// @Summary 详情
// @Description 字典详情
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.Dict} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dict/{id} [get]
func (s *DictController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Dict
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
// @Tags 字典管理
// @Summary 创建
// @Description 字典创建
// @Param token header string true "认证Token"
// @Param data body request.DictCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.Dict} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dict [post]
func (s *DictController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Dict
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
// @Tags 字典管理
// @Summary 更新
// @Description 字典更新
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Param data body request.DictUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dict/{id} [put]
func (s *DictController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.Dict
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

// Delete 删除
// @Tags 字典管理
// @Summary 删除
// @Description 字典删除
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dict/{id} [delete]
func (s *DictController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.Dict
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
