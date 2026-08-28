package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/grpc/proto"
	grpcrequest "gin/grpc/request"
	"gin/pkg/serviceprovider/debugger"
	client "gin/pkg/serviceprovider/grpcclient"
	"gin/pkg/serviceprovider/message"
	"testing"
	"time"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestUserService 测试用户grpc服务
func TestUserService(t *testing.T) {
	const traceId = "test-grpc-trace"

	facade.Register[*debugger.Debugger]("debugger", debugger.NewDebugger(message.NewEvent()))
	facade.Debugger().Start()
	defer facade.Debugger().Stop()
	defer debugger.Store.Delete(traceId)

	srv, err := client.NewServer("127.0.0.1", 1234, func(s *grpclib.Server) {
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

	ctx, cancel := context.WithTimeout(ctxkey.WithValue(context.Background(), ctxkey.TraceIdKey, traceId), 5*time.Second)
	defer cancel()

	resp, err := user.Detail(ctx, &proto.UserRequest{Id: 1})
	if err != nil {
		t.Fatalf("调用用户grpc服务失败: %v", err)
	}
	if resp.Id != 1 || resp.Username != "zhangsan" || resp.FullName != "张三" || resp.Nickname != "小张" || resp.Email != "zhangsan@example.com" || len(resp.UserRoles) != 1 || resp.UserRoles[0].Name != "管理员" {
		t.Fatalf("用户grpc响应异常: %+v", resp)
	}

	if traces := debugger.Store.Get(traceId).Grpc; len(traces) != 1 {
		t.Fatalf("grpc追踪记录数量异常: %+v", traces)
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
		"mainDept": map[string]any{
			"departmentId": 1,
		},
		"userDepts": []any{
			map[string]any{"departmentId": 1},
		},
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

// fakeUserService 用户业务测试服务
type fakeUserService struct {
	proto.UnimplementedUserServiceServer
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
