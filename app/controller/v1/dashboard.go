package v1

import (
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘
type DashboardController struct {
	base.BaseController
	service service.DashboardService
}

// Cards 卡片统计
// @Tags 仪表盘管理
// @Summary 卡片统计
// @Description 仪表盘卡片统计 用户数/角色数/菜单数/文章数/今日日志/字典项/系统配置/导入记录
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=[]service.DashboardCard} "成功"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dashboard/cards [get]
func (s *DashboardController) Cards(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := s.service.Cards(ctx)
	if err != nil {
		s.Response.Error(c, err)
		return
	}
	s.Response.Success(c, errcode.Success().WithData(res))
}

// Statistics 操作日志统计
// @Tags 仪表盘管理
// @Summary 操作日志统计
// @Description 操作日志统计 方法统计/耗时分布/小时请求量/状态码/PV/UV/平均响应/峰值时段
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=service.OperatorLogStatistics} "成功"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dashboard/statistics [get]
func (s *DashboardController) Statistics(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := s.service.Statistics(ctx)
	if err != nil {
		s.Response.Error(c, err)
		return
	}
	s.Response.Success(c, errcode.Success().WithData(res))
}

// SystemResource 系统资源
// @Tags 仪表盘管理
// @Summary 系统资源
// @Description 系统资源监控 CPU/内存/磁盘/带宽
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=service.SystemResource} "成功"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/dashboard/system-resource [get]
func (s *DashboardController) SystemResource(c *gin.Context) {
	ctx := c.Request.Context()
	res := s.service.SystemResource(ctx)
	s.Response.Success(c, errcode.Success().WithData(res))
}
