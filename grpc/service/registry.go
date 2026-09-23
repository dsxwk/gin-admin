package service

import (
	"gin/pkg/serviceprovider/grpcclient"
)

// All gRPC服务列表
func All() []grpcclient.Service {
	return []grpcclient.Service{
		UserService{},
	}
}
