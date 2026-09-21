package service

import (
	"gin/pkg/serviceprovider/grpcclient"
)

// Services 获取gRPC服务列表
func Services() []grpcclient.Service {
	return []grpcclient.Service{
		UserService{},
	}
}
