package service

import (
	"context"
	"errors"
	"gin/app/request"
	"gin/common/base"
	"gin/grpc/model"
	"gin/grpc/proto"
	grpcrequest "gin/grpc/request"
	"gin/pkg"
	client "gin/pkg/serviceprovider/grpcclient"
	"time"

	"github.com/go-viper/mapstructure/v2"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
)

func init() {
	client.Register(func(s *grpclib.Server) {
		proto.RegisterUserServiceServer(s, UserService{})
	})
	client.RegisterAuth("grpc.UserService", UserService{})
}

// UserService 用户grpc服务
type UserService struct {
	base.BaseService
	proto.UnimplementedUserServiceServer
}

// AuthMethods 方法鉴权配置
func (UserService) AuthMethods() map[string]bool {
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
func (s UserService) List(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	dto := toUserListRequest(req)
	if err := dto.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.WithContext(ctx)

	var m []model.User
	db := s.DB(&model.User{}).Model(&m)

	db = s.Search(db, m, dto.Search).
		Model(&m).
		Preload("UserRoles").
		Preload("MainDept", "is_main = ?", 1).
		Preload("MainDept.Dept").
		Preload("UserDepts").
		Preload("UserDepts.Dept")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if dto.NotPage {
		if err := db.Order("id DESC").Find(&m).Error; err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		offset, limit := request.Pagination(dto.Page, dto.PageSize)
		if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&m).Error; err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	users := make([]*proto.User, len(m))
	for i := range m {
		users[i] = toProtoUser(&m[i])
	}
	return &proto.UserListResponse{
		Total:    int32(total),
		Page:     int32(dto.Page),
		PageSize: int32(dto.PageSize),
		List:     users,
	}, nil
}

// Detail 用户详情
func (s UserService) Detail(ctx context.Context, req *proto.UserRequest) (*proto.User, error) {
	dto := grpcrequest.UserRequest{Id: req.GetId()}
	if err := dto.Validate(dto, "Detail"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.WithContext(ctx)

	var m model.User
	db := s.DB(&model.User{}).Model(&m).
		Preload("UserRoles").
		Preload("MainDept", "is_main = ?", 1).
		Preload("MainDept.Dept").
		Preload("UserDepts").
		Preload("UserDepts.Dept").
		First(&m, dto.Id)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoUser(&m), nil
}

// Create 创建用户
func (s UserService) Create(ctx context.Context, req *proto.UserCreateRequest) (*proto.User, error) {
	dto := toUserCreateRequest(req)
	if err := dto.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.WithContext(ctx)

	db := s.DB(&model.User{})

	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", dto.Username).Count(&count).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if count > 0 {
		return nil, status.Error(codes.Internal, "用户名已存在")
	}

	m := model.User{
		Username: dto.Username,
		FullName: dto.FullName,
		Nickname: dto.Nickname,
		Gender:   dto.Gender,
		Age:      dto.Age,
		Password: pkg.BcryptHash(dto.Password),
	}

	tx := db.Begin()

	if err := tx.Model(&model.User{}).Create(&m).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	if dto.UserRoles != nil {
		var userRoles []model.UserRoles
		for _, v := range dto.UserRoles {
			userRoles = append(userRoles, model.UserRoles{
				UserId: m.ID,
				RoleId: v.RoleId,
				Name:   v.Name,
			})
		}
		if err := tx.Model(&model.UserRoles{}).Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if len(dto.UserDepts) > 0 {
		var newUserDepts []model.UserDepartments
		for _, v := range dto.UserDepts {
			newUserDepts = append(newUserDepts, model.UserDepartments{
				UserId:       m.ID,
				DepartmentId: v.DepartmentId,
			})
		}
		if err := tx.Model(&model.UserDepartments{}).Create(&newUserDepts).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}

		if dto.MainDept.DepartmentId > 0 {
			mainDeptId := dto.MainDept.DepartmentId
			found := false
			for _, v := range dto.UserDepts {
				if v.DepartmentId == mainDeptId {
					found = true
					break
				}
			}
			if !found {
				tx.Rollback()
				return nil, status.Error(codes.Internal, "主部门必须为所属部门其中一个")
			}
			if err := tx.Model(&model.UserDepartments{}).
				Where("user_id = ? AND department_id = ?", m.ID, mainDeptId).
				Update("is_main", 1).Error; err != nil {
				tx.Rollback()
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoUser(&m), nil
}

// Update 更新用户
func (s UserService) Update(ctx context.Context, req *proto.UserUpdateRequest) (*proto.User, error) {
	if req.GetData() == nil {
		return nil, status.Error(codes.InvalidArgument, "更新数据不能为空")
	}
	s.WithContext(ctx)

	data := structToMap(req.GetData())
	var dto grpcrequest.UserUpdateRequest
	if err := mapstructure.Decode(data, &dto); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	dto.Id = req.GetId()
	if err := dto.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	db := s.DB(&model.User{})

	var count int64
	if err := db.Model(&model.User{}).Where("username = ? AND id <> ?", dto.Username, dto.Id).Count(&count).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if count > 0 {
		return nil, status.Error(codes.Internal, "用户名已存在")
	}

	if v, ok := data["password"].(string); ok && v != "" {
		data["password"] = pkg.BcryptHash(v)
	}
	rows := model.FilterFields(db, model.User{}, data)
	rows[model.UpdatedField] = time.Now()

	tx := db.Begin()

	if err := tx.Model(&model.User{}).Where("id = ?", dto.Id).Updates(rows).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(dto.UserRoles) > 0 {
		if err := tx.Model(&model.UserRoles{}).Where("user_id = ?", dto.Id).Delete(&model.UserRoles{}).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}

		var newUserRoles []model.UserRoles
		for _, item := range dto.UserRoles {
			newUserRoles = append(newUserRoles, model.UserRoles{
				UserId: dto.Id,
				RoleId: item.RoleId,
				Name:   item.Name,
			})
		}
		if err := tx.Model(&model.UserRoles{}).Create(&newUserRoles).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if len(dto.UserDepts) > 0 {
		if err := tx.Model(&model.UserDepartments{}).Where("user_id = ?", dto.Id).Delete(&model.UserDepartments{}).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}

		var newUserDepts []model.UserDepartments
		for _, item := range dto.UserDepts {
			newUserDepts = append(newUserDepts, model.UserDepartments{
				UserId:       dto.Id,
				DepartmentId: item.DepartmentId,
			})
		}
		if err := tx.Model(&model.UserDepartments{}).Create(&newUserDepts).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, err.Error())
		}

		if dto.MainDept.DepartmentId > 0 {
			mainDeptId := dto.MainDept.DepartmentId
			found := false
			for _, item := range dto.UserDepts {
				if item.DepartmentId == mainDeptId {
					found = true
					break
				}
			}
			if !found {
				tx.Rollback()
				return nil, status.Error(codes.Internal, "主部门必须为所属部门其中一个")
			}
			if err := tx.Model(&model.UserDepartments{}).
				Where("user_id = ? AND department_id = ?", dto.Id, mainDeptId).
				Update("is_main", 1).Error; err != nil {
				tx.Rollback()
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return s.Detail(ctx, &proto.UserRequest{Id: dto.Id})
}

// Delete 删除用户
func (s UserService) Delete(ctx context.Context, req *proto.UserRequest) (*proto.EmptyResponse, error) {
	dto := grpcrequest.UserRequest{Id: req.GetId()}
	if err := dto.Validate(dto, "Delete"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.WithContext(ctx)

	var m model.User
	if err := s.DB(&model.User{}).Model(&m).Delete(&m, dto.Id).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.EmptyResponse{}, nil
}

// BatchDelete 批量删除用户
func (s UserService) BatchDelete(ctx context.Context, req *proto.UserBatchDeleteRequest) (*proto.EmptyResponse, error) {
	dto := grpcrequest.UserBatchDeleteRequest{Ids: req.GetIds()}
	if err := dto.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.WithContext(ctx)

	var m model.User
	if err := s.DB(&model.User{}).Model(&m).Delete(&m, dto.Ids).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.EmptyResponse{}, nil
}

// toUserListRequest 列表请求转换
func toUserListRequest(req *proto.UserListRequest) *grpcrequest.UserListRequest {
	if req == nil {
		return &grpcrequest.UserListRequest{}
	}
	return &grpcrequest.UserListRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		NotPage:  req.NotPage,
		Search:   structToMap(req.Search),
		Sort:     structToMap(req.Sort),
	}
}

// toUserCreateRequest 创建请求转换
func toUserCreateRequest(req *proto.UserCreateRequest) *grpcrequest.UserCreateRequest {
	if req == nil {
		return &grpcrequest.UserCreateRequest{}
	}
	dto := &grpcrequest.UserCreateRequest{
		Username: req.Username,
		FullName: req.FullName,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Password: req.Password,
		Age:      req.Age,
		MainDept: grpcrequest.DeptItem{DepartmentId: req.GetMainDept().GetDepartmentId()},
	}
	for _, item := range req.UserRoles {
		dto.UserRoles = append(dto.UserRoles, grpcrequest.UserRoleItem{
			RoleId: item.GetRoleId(),
			Name:   item.GetName(),
		})
	}
	for _, item := range req.UserDepts {
		dto.UserDepts = append(dto.UserDepts, grpcrequest.DeptItem{DepartmentId: item.GetDepartmentId()})
	}
	return dto
}

// structToMap 结构转map
func structToMap(s *structpb.Struct) map[string]any {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// toProtoUser 用户模型转proto用户
func toProtoUser(m *model.User) *proto.User {
	if m == nil {
		return nil
	}
	return &proto.User{
		Id:        m.ID,
		Avatar:    m.Avatar,
		Username:  m.Username,
		FullName:  m.FullName,
		Email:     m.Email,
		Nickname:  m.Nickname,
		Gender:    m.Gender,
		Age:       m.Age,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		UserRoles: toProtoUserRoles(m.UserRoles),
		MainDept:  toProtoUserDepartment(m.MainDept),
		UserDepts: toProtoUserDepartments(m.UserDepts),
	}
}

// toProtoUserRoles 用户角色转换
func toProtoUserRoles(items []*model.UserRoles) []*proto.UserRole {
	if len(items) == 0 {
		return nil
	}
	result := make([]*proto.UserRole, len(items))
	for i, item := range items {
		result[i] = &proto.UserRole{
			Id:        item.ID,
			UserId:    item.UserId,
			RoleId:    item.RoleId,
			Name:      item.Name,
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
	}
	return result
}

// toProtoUserDepartment 用户部门转换
func toProtoUserDepartment(item *model.UserDepartments) *proto.UserDepartment {
	if item == nil {
		return nil
	}
	return &proto.UserDepartment{
		Id:           item.ID,
		UserId:       item.UserId,
		DepartmentId: item.DepartmentId,
		IsMain:       item.IsMain,
		CreatedAt:    item.CreatedAt.String(),
		UpdatedAt:    item.UpdatedAt.String(),
	}
}

// toProtoUserDepartments 用户部门列表转换
func toProtoUserDepartments(items []*model.UserDepartments) []*proto.UserDepartment {
	if len(items) == 0 {
		return nil
	}
	result := make([]*proto.UserDepartment, len(items))
	for i, item := range items {
		result[i] = toProtoUserDepartment(item)
	}
	return result
}
