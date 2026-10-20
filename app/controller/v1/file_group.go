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

// FileGroupController 文件分组控制器
type FileGroupController struct {
	base.Controller
	service service.FileGroupService
}

// List 文件分组列表
// @Tags 文件分组管理
// @Summary 文件分组列表
// @Description 查询文件分组列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.FileGroup}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file-group [get]
func (s *FileGroupController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.FileGroup
	)

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

// Create 添加文件分组
// @Tags 文件分组管理
// @Summary 添加文件分组
// @Description 添加文件分组
// @Param token header string true "认证Token"
// @Param data body request.FileGroup true "添加参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.FileGroup} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file-group [post]
func (s *FileGroupController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.FileGroup
	)

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

// Update 修改文件分组
// @Tags 文件分组管理
// @Summary 修改文件分组
// @Description 修改文件分组
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Param data body request.FileGroup true "修改参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file-group/{id} [put]
func (s *FileGroupController) Update(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		data map[string]any
		req  request.FileGroup
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

	err = s.service.Update(ctx, req.ID, req, data)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(data))
}

// Delete 删除文件分组
// @Tags 文件分组管理
// @Summary 删除文件分组
// @Description 删除文件分组并将文件移回未分组
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file-group/{id} [delete]
func (s *FileGroupController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.FileGroup
	)

	req.ID = facade.Request().Path[int64](c, "id", 0)
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
