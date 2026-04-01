package main

import (
	"gin/bootstrap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go get # -u
//go:generate go mod tidy
//go:generate go mod download
//go:generate go mod vendor

// @title Gin Swagger API
// @version 2.0
// @description Gin API 文档
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email 25076778@qq.com
// @host 127.0.0.1:8080
func main() {
	// 初始化应用
	bootstrap.NewApp().Run()
}
