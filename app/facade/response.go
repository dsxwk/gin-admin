package facade

import "gin/pkg/errcode"

// Response 响应门面方法
// 使用示例:
//
//	facade.Response().Success(c, errcode.Success().WithData(data))
//	facade.Response().Error(c, err)
func Response() errcode.Response {
	return errcode.Response{}
}
