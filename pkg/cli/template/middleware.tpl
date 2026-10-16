package {{.Package}}

import (
    "github.com/gin-gonic/gin"
    "gin/common/base"
)

type {{.Name}} struct {
	base.Middleware
}

// Handle {{.Description}}
func (s {{.Name}}) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Define your middleware logic here
        c.Next()
    }
}
