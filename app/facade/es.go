package facade

import (
	"gin/pkg/serviceprovider/es"
)

// ES ES客户端门面
// 使用示例:
//
//	client := facade.ES()
func ES() *es.Client {
	return Get[*es.Client]("es")
}
