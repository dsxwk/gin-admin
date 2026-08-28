package facade

import (
	"fmt"
	"gin/pkg/serviceprovider/grpcclient"

	grpclib "google.golang.org/grpc"
)

// Grpc 门面函数
// 使用示例:
//
//	user, err := facade.Grpc().Service(proto.NewUserServiceClient)
//	resp, err := user.Detail(ctx, req)
func Grpc() GrpcFacade {
	return GrpcFacade{client: Get[*grpcclient.Client]("grpc")}
}

type GrpcFacade struct {
	client *grpcclient.Client
}

// Client 获取grpc客户端
func (g GrpcFacade) Client() *grpcclient.Client {
	return g.client
}

// Service 获取指定服务客户端
func (g GrpcFacade) Service[T any](factory func(grpclib.ClientConnInterface) T) (T, error) {
	var zero T
	if g.client == nil {
		return zero, fmt.Errorf("grpc客户端未初始化")
	}
	return g.client.Service(factory), nil
}
