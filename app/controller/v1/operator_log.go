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

type OperatorLogController struct {
	base.BaseController
	service service.OperatorLogService
}

// List 列表
// @Tags 操作日志管理
// @Summary 列表
// @Description 操作日志列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Param __search query string false "搜索条件(JSON)"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.OperatorLog}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/operator-log [get]
func (s *OperatorLogController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.OperatorLog
	)

	s.service.WithContext(ctx)

	err := facade.Request().BindValidate(c, &req, "List")
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
// @Tags 操作日志管理
// @Summary 详情
// @Description 操作日志详情
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.OperatorLog} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/operator-log/{id} [get]
func (s *OperatorLogController) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.OperatorLog
	)

	s.service.WithContext(ctx)

	req.ID = facade.Request().Path[int64](c, "id", 0)

	err := facade.Request().BindValidate(c, &req, "Detail")
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

// Delete 删除
// @Tags 操作日志管理
// @Summary 删除
// @Description 操作日志删除
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/operator-log/{id} [delete]
func (s *OperatorLogController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.OperatorLog
	)

	s.service.WithContext(ctx)

	req.ID = facade.Request().Path[int64](c, "id", 0)

	err := facade.Request().BindValidate(c, &req, "Delete")
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

// BatchDelete 批量删除
// @Tags 操作日志管理
// @Summary 批量删除
// @Description 操作日志批量删除
// @Param token header string true "认证Token"
// @Param data body request.OperatorLog true "批量删除参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/operator-log/batch-delete [post]
func (s *OperatorLogController) BatchDelete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.OperatorLog
	)

	s.service.WithContext(ctx)

	err := facade.Request().BindValidate(c, &req, "BatchDelete")
	if err != nil {
		s.Response.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = s.service.BatchDelete(req.IDs)
	if err != nil {
		s.Response.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Response.Success(c, errcode.Success())
}
