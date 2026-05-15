package {{.Package}}

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type {{.Name}}Controller struct {
    base.BaseController
}

// {{.Function}} {{.Description}}
// @Router {{.Router}} [{{.Method}}]
func (s *{{.Name}}Controller) {{.Function}}(c *gin.Context) {
    // Define your function here
    s.Response.Success(c, errcode.Success().WithMsg("Test Msg").WithData([]string{}))
}