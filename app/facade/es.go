package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/es"
)

// ES ES客户端门面
// 使用示例:
//
//	client := facade.ES()
func ES() *es.Client {
	return container.Default().Get[*es.Client](serviceprovider.ServiceES)
}
