package base

import "gin/common/response"

type BaseMiddleware struct {
	Context
	Response response.Response
}
