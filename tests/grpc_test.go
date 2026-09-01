package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/grpc/proto"
	grpcrequest "gin/grpc/request"
	client "gin/pkg/serviceprovider/grpcclient"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	grpcTestJwtKey   = "test-jwt-key"
	grpcTestPort     = 1234
	grpcAuthTestPort = 1235
)

// TestGrpcMethodAuth 测试grpc方法级JWT鉴权
func TestGrpcMethodAuth(t *testing.T) {
	client.SetJwtKey(grpcTestJwtKey)
	client.RegisterAuth("grpc.AuthTestService", authTestServer{})

	srv, err := client.NewServer("127.0.0.1", grpcAuthTestPort, func(s *grpclib.Server) {
		s.RegisterService(&grpclib.ServiceDesc{
			ServiceName: "grpc.AuthTestService",
			HandlerType: (*authTestService)(nil),
			Methods: []grpclib.MethodDesc{
				{
					MethodName: "Ping",
					Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpclib.UnaryServerInterceptor) (any, error) {
						in := new(emptypb.Empty)
						if err := dec(in); err != nil {
							return nil, err
						}
						if interceptor == nil {
							return srv.(authTestServer).Ping(ctx, in)
						}
						info := &grpclib.UnaryServerInfo{Server: srv, FullMethod: "/grpc.AuthTestService/Ping"}
						handler := func(ctx context.Context, req any) (any, error) {
							return srv.(authTestServer).Ping(ctx, req.(*emptypb.Empty))
						}
						return interceptor(ctx, in, info, handler)
					},
				},
				{
					MethodName: "PingPublic",
					Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpclib.UnaryServerInterceptor) (any, error) {
						in := new(emptypb.Empty)
						if err := dec(in); err != nil {
							return nil, err
						}
						if interceptor == nil {
							return srv.(authTestServer).PingPublic(ctx, in)
						}
						info := &grpclib.UnaryServerInfo{Server: srv, FullMethod: "/grpc.AuthTestService/PingPublic"}
						handler := func(ctx context.Context, req any) (any, error) {
							return srv.(authTestServer).PingPublic(ctx, req.(*emptypb.Empty))
						}
						return interceptor(ctx, in, info, handler)
					},
				},
			},
		}, authTestServer{})
	})
	if err != nil {
		t.Fatalf("创建grpc鉴权服务端失败: %v", err)
	}
	if err = srv.Start(); err != nil {
		t.Fatalf("启动grpc鉴权服务端失败: %v", err)
	}
	defer srv.Stop()

	c, err := client.NewClient(srv.Addr())
	if err != nil {
		t.Fatalf("创建grpc鉴权客户端失败: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = c.Conn.Invoke(ctx, "/grpc.AuthTestService/Ping", &emptypb.Empty{}, &emptypb.Empty{}); err == nil || status.Code(err) != codes.Unauthenticated {
		t.Fatalf("未携带Token应返回未鉴权: %v", err)
	}

	if err = c.Conn.Invoke(ctx, "/grpc.AuthTestService/PingPublic", &emptypb.Empty{}, &emptypb.Empty{}); err != nil {
		t.Fatalf("免鉴权方法调用失败: %v", err)
	}

	if err = c.Conn.Invoke(client.WithToken(ctx, "invalid"), "/grpc.AuthTestService/Ping", &emptypb.Empty{}, &emptypb.Empty{}); err == nil || status.Code(err) != codes.Unauthenticated {
		t.Fatalf("无效Token应返回未鉴权: %v", err)
	}

	if err = c.Conn.Invoke(client.WithToken(ctx, newGrpcTestToken(1)), "/grpc.AuthTestService/Ping", &emptypb.Empty{}, &emptypb.Empty{}); err != nil {
		t.Fatalf("有效Token调用失败: %v", err)
	}
}

// TestGrpcUserService 测试用户grpc服务
func TestGrpcUserService(t *testing.T) {
	client.SetJwtKey(grpcTestJwtKey)
	client.RegisterAuth("grpc.UserService", fakeUserService{})

	srv, err := client.NewServer("127.0.0.1", grpcTestPort, func(s *grpclib.Server) {
		proto.RegisterUserServiceServer(s, fakeUserService{})
	})
	if err != nil {
		t.Fatalf("创建grpc服务端失败: %v", err)
	}
	if err = srv.Start(); err != nil {
		t.Fatalf("启动grpc服务端失败: %v", err)
	}
	defer srv.Stop()

	c, err := client.NewClient(srv.Addr())
	if err != nil {
		t.Fatalf("创建grpc客户端失败: %v", err)
	}
	defer c.Close()

	facade.Register[*client.Client]("grpc", c)

	user, err := facade.Grpc().Service(proto.NewUserServiceClient)
	if err != nil {
		t.Fatalf("获取用户grpc客户端失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = facade.Grpc().WithToken(ctx, newGrpcTestToken(1))

	resp, err := user.Detail(ctx, &proto.UserRequest{Id: 1})
	if err != nil {
		t.Fatalf("调用用户grpc服务失败: %v", err)
	}
	if resp.Id != 1 || resp.Username != "zhangsan" || resp.FullName != "张三" || resp.Nickname != "小张" || resp.Email != "zhangsan@example.com" || len(resp.UserRoles) != 1 || resp.UserRoles[0].Name != "管理员" {
		t.Fatalf("用户grpc响应异常: %+v", resp)
	}

	_, err = user.Detail(ctx, &proto.UserRequest{})
	if err == nil || status.Code(err) != codes.InvalidArgument {
		t.Fatalf("用户grpc验证异常: %v", err)
	}

	_, err = user.Detail(ctx, &proto.UserRequest{Id: 99})
	if err == nil || status.Code(err) != codes.NotFound {
		t.Fatalf("用户grpc不存在异常: %v", err)
	}

	listResp, err := user.List(ctx, &proto.UserListRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("调用用户grpc列表失败: %v", err)
	}
	if listResp.Total != 1 || len(listResp.List) != 1 || listResp.List[0].Username != "zhangsan" {
		t.Fatalf("用户grpc列表响应异常: %+v", listResp)
	}

	createResp, err := user.Create(ctx, &proto.UserCreateRequest{
		Username: "lisi",
		FullName: "李四",
		Nickname: "小李",
		Gender:   1,
		Password: "123456",
		Age:      20,
		MainDept: &proto.DeptItem{DepartmentId: 1},
		UserDepts: []*proto.DeptItem{
			{DepartmentId: 1},
		},
	})
	if err != nil {
		t.Fatalf("调用用户grpc创建失败: %v", err)
	}
	if createResp.Username != "lisi" || createResp.FullName != "李四" {
		t.Fatalf("用户grpc创建响应异常: %+v", createResp)
	}

	updateData, _ := structpb.NewStruct(map[string]any{
		"username": "lisi",
		"fullName": "李四",
		"nickname": "李四四",
		"gender":   1,
		"age":      21,
	})
	updateResp, err := user.Update(ctx, &proto.UserUpdateRequest{
		Id:   2,
		Data: updateData,
	})
	if err != nil {
		t.Fatalf("调用用户grpc更新失败: %v", err)
	}
	if updateResp.Nickname != "李四四" {
		t.Fatalf("用户grpc更新响应异常: %+v", updateResp)
	}

	if _, err = user.Delete(ctx, &proto.UserRequest{Id: 2}); err != nil {
		t.Fatalf("调用用户grpc删除失败: %v", err)
	}
	if _, err = user.BatchDelete(ctx, &proto.UserBatchDeleteRequest{
		Ids: []int32{2, 3},
	}); err != nil {
		t.Fatalf("调用用户grpc批量删除失败: %v", err)
	}

}

// newGrpcTestToken 生成测试JWT
func newGrpcTestToken(id int64) string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  float64(id),
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(grpcTestJwtKey))
	if err != nil {
		panic(err)
	}
	return token
}

// authTestServer 鉴权测试服务
type authTestServer struct{}

// authTestService 鉴权测试服务接口
type authTestService interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	PingPublic(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

// AuthMethods 方法鉴权配置
func (authTestServer) AuthMethods() map[string]bool {
	return map[string]bool{
		"Ping":       true,
		"PingPublic": false,
	}
}

// Ping 测试方法
func (authTestServer) Ping(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	if ctxkey.GetValue(ctx, ctxkey.UserIdKey) != int64(1) {
		return nil, status.Error(codes.Unauthenticated, "用户ID缺失")
	}
	return &emptypb.Empty{}, nil
}

// PingPublic 免鉴权测试方法
func (authTestServer) PingPublic(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// fakeUserService 用户业务测试服务
type fakeUserService struct {
	proto.UnimplementedUserServiceServer
}

// AuthMethods 方法鉴权配置
func (fakeUserService) AuthMethods() map[string]bool {
	return map[string]bool{
		"Detail":      true,
		"List":        true,
		"Create":      true,
		"Update":      true,
		"Delete":      true,
		"BatchDelete": true,
	}
}

// List 用户列表
func (h fakeUserService) List(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	return &proto.UserListResponse{
		Total:    1,
		Page:     req.Page,
		PageSize: req.PageSize,
		List: []*proto.User{{
			Id:        1,
			Username:  "zhangsan",
			FullName:  "张三",
			Nickname:  "小张",
			CreatedAt: "2026-01-01 00:00:00",
			UpdatedAt: "2026-01-01 00:00:00",
		}},
	}, nil
}

// Detail 用户详情
func (h fakeUserService) Detail(ctx context.Context, req *proto.UserRequest) (*proto.User, error) {
	dto := grpcrequest.UserRequest{Id: req.GetId()}
	if err := dto.Validate(dto, "Detail"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.GetId() != 1 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return &proto.User{
		Id:        1,
		Avatar:    "avatar.png",
		Username:  "zhangsan",
		FullName:  "张三",
		Email:     "zhangsan@example.com",
		Nickname:  "小张",
		Gender:    1,
		Age:       18,
		Status:    1,
		CreatedAt: "2026-01-01 00:00:00",
		UpdatedAt: "2026-01-01 00:00:00",
		UserRoles: []*proto.UserRole{{
			Id:        1,
			UserId:    1,
			RoleId:    1,
			Name:      "管理员",
			CreatedAt: "2026-01-01 00:00:00",
			UpdatedAt: "2026-01-01 00:00:00",
		}},
		MainDept: &proto.UserDepartment{
			Id:           1,
			UserId:       1,
			DepartmentId: 1,
			IsMain:       1,
			CreatedAt:    "2026-01-01 00:00:00",
			UpdatedAt:    "2026-01-01 00:00:00",
		},
		UserDepts: []*proto.UserDepartment{{
			Id:           1,
			UserId:       1,
			DepartmentId: 1,
			IsMain:       1,
			CreatedAt:    "2026-01-01 00:00:00",
			UpdatedAt:    "2026-01-01 00:00:00",
		}},
	}, nil
}

// Create 创建用户
func (h fakeUserService) Create(ctx context.Context, req *proto.UserCreateRequest) (*proto.User, error) {
	return &proto.User{
		Id:        2,
		Username:  req.Username,
		FullName:  req.FullName,
		Nickname:  req.Nickname,
		Gender:    req.Gender,
		Age:       req.Age,
		CreatedAt: "2026-01-01 00:00:00",
		UpdatedAt: "2026-01-01 00:00:00",
	}, nil
}

// Update 更新用户
func (h fakeUserService) Update(ctx context.Context, req *proto.UserUpdateRequest) (*proto.User, error) {
	data := req.GetData().AsMap()
	return &proto.User{
		Id:        req.Id,
		Username:  data["username"].(string),
		FullName:  data["fullName"].(string),
		Nickname:  data["nickname"].(string),
		Gender:    int32(data["gender"].(float64)),
		Age:       int32(data["age"].(float64)),
		CreatedAt: "2026-01-01 00:00:00",
		UpdatedAt: "2026-01-01 00:00:00",
	}, nil
}

// Delete 删除用户
func (h fakeUserService) Delete(ctx context.Context, req *proto.UserRequest) (*proto.EmptyResponse, error) {
	return &proto.EmptyResponse{}, nil
}

// BatchDelete 批量删除用户
func (h fakeUserService) BatchDelete(ctx context.Context, req *proto.UserBatchDeleteRequest) (*proto.EmptyResponse, error) {
	return &proto.EmptyResponse{}, nil
}
