package {{.Package}}

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type {{.Name}}Controller struct {
    base.BaseController
}

// List 列表
// @Tags {{.Description}}管理
// @Summary 列表
// @Description {{.Description}}列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Param notPage query string true "是否不分页"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.{{.Name}}}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router {{.RoutePath}} [get]
func (s *{{.Name}}Controller) List(c *gin.Context) {
	// todo
	s.Response.Success(c, errcode.Success())
}

// Detail 详情
// @Tags {{.Description}}管理
// @Summary 详情
// @Description {{.Description}}详情
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.{{.Name}}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router {{.RoutePath}}/{id} [get]
func (s *{{.Name}}Controller) Detail(c *gin.Context) {
	// todo
	s.Response.Success(c, errcode.Success())
}

// Create 创建
// @Tags {{.Description}}管理
// @Summary 创建
// @Description {{.Description}}创建
// @Param token header string true "认证Token"
// @Param data body request.RoleCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.{{.Name}}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router {{.RoutePath}} [post]
func (s *{{.Name}}Controller) Create(c *gin.Context) {
	// todo
	s.Response.Success(c, errcode.Success())
}

// Update 更新
// @Tags {{.Description}}管理
// @Summary 更新
// @Description {{.Description}}更新
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Param data body request.RoleUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router {{.RoutePath}}/{id} [put]
func (s *{{.Name}}Controller) Update(c *gin.Context) {
	// todo
	s.Response.Success(c, errcode.Success())
}

// Delete 删除
// @Tags {{.Description}}管理
// @Summary 删除
// @Description {{.Description}}删除
// @Param token header string true "认证Token"
// @Param id path int true "ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router {{.RoutePath}}/{id} [delete]
func (s *{{.Name}}Controller) Delete(c *gin.Context) {
	// todo
	s.Response.Success(c, errcode.Success())
}
