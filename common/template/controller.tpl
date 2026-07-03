package {{.Package}}

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

type {{.Name}}Controller struct {
    base.BaseController
    service service.{{.Name}}Service
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
	var (
        ctx = c.Request.Context()
        req request.{{.Name}}
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
	var (
        ctx = c.Request.Context()
        req request.{{.Name}}
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
	var (
        ctx = c.Request.Context()
        req request.{{.Name}}
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
	var (
        ctx  = c.Request.Context()
        data map[string]interface{}
        req  request.{{.Name}}
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
	var (
        ctx = c.Request.Context()
        req request.{{.Name}}
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
