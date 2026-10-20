package v1

import (
	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	base.Controller
	service service.FileService
}

// List 文件列表
// @Tags 文件管理
// @Summary 文件列表
// @Description 查询文件列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Param groupId query int false "分组ID"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.File}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file [get]
func (s *FileController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.File
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

// Upload 上传文件
// @Tags 文件管理
// @Summary 上传文件
// @Description 上传文件并保存文件记录
// @Accept multipart/form-data
// @Produce json
// @Param token header string true "认证Token"
// @Param file formData file true "上传文件"
// @Param groupId formData int false "分组ID"
// @Param storage formData string false "存储磁盘"
// @Success 200 {object} errcode.SuccessResponse{data=model.File} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file/upload [post]
func (s *FileController) Upload(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg("请上传file文件"))
		return
	}

	m, err := s.service.Upload(
		ctx,
		s.GetUserID(c),
		facade.Request().Form[int64](c, "groupId", 0),
		facade.Request().Form[string](c, "storage", ""),
		fileHeader,
	)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}

// Delete 删除文件
// @Tags 文件管理
// @Summary 删除文件
// @Description 删除文件记录和物理文件
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file/{id} [delete]
func (s *FileController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.File
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

// UpdateGroup 修改文件分组
// @Tags 文件管理
// @Summary 修改文件分组
// @Description 修改文件所属分组,groupId为0表示移回未分组
// @Param token header string true "认证Token"
// @Param id path int true "文件ID"
// @Param data body request.FileUpdateGroup true "分组参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.File} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/file/{id}/group [put]
func (s *FileController) UpdateGroup(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.FileUpdateGroup
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}
	if req.GroupId == nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg("groupId必填"))
		return
	}

	req.ID = facade.Request().Path[int64](c, "id", 0)
	err = facade.Request().Validate(c, &req, "UpdateGroup")
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	m, err := s.service.UpdateGroup(ctx, req.ID, *req.GroupId)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(m))
}
