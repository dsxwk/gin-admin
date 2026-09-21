package grpcclient

import (
	grpclib "google.golang.org/grpc"
)

// Service grpc服务接口
type Service interface {
	Name() string
	Register(s *grpclib.Server)
}
