package v1

import (
	"gin/app/model"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/errcode"
	"gin/pkg/lang"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"strconv"
)

type UserController struct {
	base.BaseController
	service service.UserService
	req     request.User
	search  request.Search
	user    model.User
}

// List 列表
// @Tags 用户管理
// @Summary 列表
// @Description 用户列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.User}} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user [get]
func (s *UserController) List(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)

	s.service.WithContext(ctx)

	err := c.ShouldBind(&s.search)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	err = c.ShouldBind(&s.req)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	// 验证
	err = s.req.Validate(s.req, "List")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	res, err := s.service.List(s.req, s.search.Search)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(lang.T(ctx, err.Error(), nil)))
		return
	}

	s.Success(c, errcode.Success().WithData(res))
}

// Create 创建
// @Tags 用户管理
// @Summary 创建
// @Description 用户创建
// @Param token header string true "认证Token"
// @Param data body request.UserCreate true "创建参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user [post]
func (s *UserController) Create(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)

	s.service.WithContext(ctx)

	err := c.ShouldBind(&s.req)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	// 验证
	err = s.req.Validate(s.req, "Create")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = copier.Copy(&s.user, &s.req)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(lang.T(ctx, err.Error(), nil)))
		return
	}

	s.user, err = s.service.Create(s.user)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(lang.T(ctx, err.Error(), nil)))
		return
	}

	s.Success(c, errcode.Success().WithData(s.user))
}

// Update 更新
// @Tags 用户管理
// @Summary 更新
// @Description 用户更新
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Param data body request.UserUpdate true "更新参数"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [put]
func (s *UserController) Update(c *gin.Context) {
	var (
		id   int64
		ctx  = c.Request.Context()
		data map[string]interface{}
	)

	s.service.WithContext(ctx)

	err := c.ShouldBindBodyWith(&s.req, binding.JSON)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	err = c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	id, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}
	s.req.ID = id

	// 验证
	err = s.req.Validate(s.req, "Update")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	err = copier.Copy(&s.user, &s.req)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	err = s.service.Update(s.user.ID, data)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Success(c, errcode.Success().WithData(data))
}

// Detail 详情
// @Tags 用户管理
// @Summary 详情
// @Description 用户详情
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [get]
func (s *UserController) Detail(c *gin.Context) {
	var (
		id  int64
		ctx = c.Request.Context()
	)

	s.service.WithContext(ctx)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}
	s.req.ID = id

	// 验证
	err = s.req.Validate(s.req, "Detail")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.Detail(s.req.ID)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Success(c, errcode.Success().WithData(m))
}

// Delete 删除
// @Tags 用户管理
// @Summary 删除
// @Description 用户删除
// @Param token header string true "认证Token"
// @Param id path int true "用户ID"
// @Success 200 {object} errcode.SuccessResponse{data=model.User} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user/{id} [delete]
func (s *UserController) Delete(c *gin.Context) {
	var (
		id  int64
		ctx = c.Request.Context()
	)

	s.service.WithContext(ctx)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}
	s.req.ID = id

	// 验证
	err = s.req.Validate(s.req, "Delete")
	if err != nil {
		s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
		return
	}

	m, err := s.service.Delete(s.req.ID)
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
	}

	s.Success(c, errcode.Success().WithData(m))
}
