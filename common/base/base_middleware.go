package base

import "gin/pkg/errcode"

type BaseMiddleware struct {
	Context
	Response errcode.Response
}
