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

type ConfigCategoryController struct {
	base.BaseController
	service service.ConfigCategoryService
}

// List 列表
// @Tags 配置分类管理
// @Summary 列表
// @Description 配置分类列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.ConfigCategory}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/config-category [get]
func (s *ConfigCategoryController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.ConfigCategory
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

// Detail 详情
// @Tags 配置分类管理
// @Summary 详情
// @Description 配置分类详情
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.ConfigCategory} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/config-category/{id} [get]
func (s *ConfigCategoryController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.ConfigCategory
	)

	s.service.WithContext(ctx)

	req.ID = facade.Request[int64]().Path(c, "id", 0)

	// 绑定参数并验证
	err := facade.Request[any]().BindValidate(c, &req, "Detail")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.Detail(req.ID)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Create 创建
// @Tags 配置分类管理
// @Summary 创建
// @Description 配置分类创建
// @Param token header string true "认证Token"
// @Param data body request.ConfigCategoryCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.ConfigCategory} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/config-category [post]
func (s *ConfigCategoryController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.ConfigCategory
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

// Update 更新
// @Tags 配置分类管理
// @Summary 更新
// @Description 配置分类更新
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Param data body request.ConfigCategoryUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/config-category/{id} [put]
func (s *ConfigCategoryController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]interface{}
		req  request.ConfigCategory
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
// @Tags 配置分类管理
// @Summary 删除
// @Description 配置分类删除
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/config-category/{id} [delete]
func (s *ConfigCategoryController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.ConfigCategory
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
