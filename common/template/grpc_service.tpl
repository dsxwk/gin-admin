package {{.Package}}

import (
    "context"
    "errors"
    "gin/app/request"
    "gin/common/base"
    "gin/grpc/model"
    "gin/grpc/proto"
    grpcrequest "gin/grpc/request"
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
        proto.Register{{.Name}}ServiceServer(s, {{.Name}}Service{})
    })
    client.RegisterAuth("grpc.{{.Name}}Service", {{.Name}}Service{})
}

// {{.Name}}Service {{.Description}}grpc服务
type {{.Name}}Service struct {
    base.BaseService
    proto.Unimplemented{{.Name}}ServiceServer
}

// AuthMethods 方法鉴权配置
func ({{.Name}}Service) AuthMethods() map[string]bool {
    return map[string]bool{
        "Detail": {{.Auth}},
        "List":   {{.Auth}},
        "Create": {{.Auth}},
        "Update": {{.Auth}},
        "Delete": {{.Auth}},
    }
}

// List {{.Description}}列表
func (s {{.Name}}Service) List(ctx context.Context, req *proto.{{.Name}}ListRequest) (*proto.{{.Name}}ListResponse, error) {
    dto := to{{.Name}}ListRequest(req)
    if err := dto.Validate(); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    var m []model.{{.Name}}
    db := s.DB(ctx, &model.{{.Name}}{}).Model(&m)

    db = s.Search(db, m, dto.Search).Model(&m)

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

    list := make([]*proto.{{.Name}}, len(m))
    for i := range m {
        list[i] = toProto{{.Name}}(&m[i])
    }
    return &proto.{{.Name}}ListResponse{
        Total:    int32(total),
        Page:     int32(dto.Page),
        PageSize: int32(dto.PageSize),
        List:     list,
    }, nil
}

// Detail {{.Description}}详情
func (s {{.Name}}Service) Detail(ctx context.Context, req *proto.{{.Name}}Request) (*proto.{{.Name}}, error) {
    dto := grpcrequest.{{.Name}}Request{Id: req.GetId()}
    if err := dto.Validate(dto, "Detail"); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    var m model.{{.Name}}
    db := s.DB(ctx, &model.{{.Name}}{}).Model(&m).First(&m, dto.Id)
    if err := db.Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Error(codes.NotFound, "{{.Description}}不存在")
        }
        return nil, status.Error(codes.Internal, err.Error())
    }
    return toProto{{.Name}}(&m), nil
}

// Create {{.Description}}创建
func (s {{.Name}}Service) Create(ctx context.Context, req *proto.{{.Name}}CreateRequest) (*proto.{{.Name}}, error) {
    dto := to{{.Name}}CreateRequest(req)
    if err := dto.Validate(); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    m := model.{{.Name}}{
{{.CreateFields}}
    }
    if err := s.DB(ctx, &model.{{.Name}}{}).Model(&m).Create(&m).Error; err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return toProto{{.Name}}(&m), nil
}

// Update {{.Description}}更新
func (s {{.Name}}Service) Update(ctx context.Context, req *proto.{{.Name}}UpdateRequest) (*proto.{{.Name}}, error) {
    if req.GetData() == nil {
        return nil, status.Error(codes.InvalidArgument, "更新数据不能为空")
    }
    data := {{.Var}}StructToMap(req.GetData())
    var dto grpcrequest.{{.Name}}UpdateRequest
    if err := mapstructure.Decode(data, &dto); err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    dto.Id = req.GetId()
    if err := dto.Validate(); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    db := s.DB(ctx, &model.{{.Name}}{})
    rows := model.FilterFields(db, model.{{.Name}}{}, data)
{{.UpdateJsonFields}}
    rows[model.UpdatedField] = time.Now()
    if err := db.Model(&model.{{.Name}}{}).Where("id = ?", req.GetId()).Updates(rows).Error; err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return s.Detail(ctx, &proto.{{.Name}}Request{Id: req.GetId()})
}

// Delete {{.Description}}删除
func (s {{.Name}}Service) Delete(ctx context.Context, req *proto.{{.Name}}Request) (*proto.EmptyResponse, error) {
    dto := grpcrequest.{{.Name}}Request{Id: req.GetId()}
    if err := dto.Validate(dto, "Delete"); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    var m model.{{.Name}}
    if err := s.DB(ctx, &model.{{.Name}}{}).Model(&m).Delete(&m, dto.Id).Error; err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return &proto.EmptyResponse{}, nil
}

// to{{.Name}}ListRequest 列表请求转换
func to{{.Name}}ListRequest(req *proto.{{.Name}}ListRequest) *grpcrequest.{{.Name}}ListRequest {
    if req == nil {
        return &grpcrequest.{{.Name}}ListRequest{}
    }
    return &grpcrequest.{{.Name}}ListRequest{
        Page:     int(req.Page),
        PageSize: int(req.PageSize),
        NotPage:  req.NotPage,
        Search:   {{.Var}}StructToMap(req.Search),
        Sort:     {{.Var}}StructToMap(req.Sort),
    }
}

// to{{.Name}}CreateRequest 创建请求转换
func to{{.Name}}CreateRequest(req *proto.{{.Name}}CreateRequest) *grpcrequest.{{.Name}}CreateRequest {
    if req == nil {
        return &grpcrequest.{{.Name}}CreateRequest{}
    }
    return &grpcrequest.{{.Name}}CreateRequest{
{{.RequestCreateFields}}
    }
}

// {{.Var}}StructToMap 结构转map
func {{.Var}}StructToMap(s *structpb.Struct) map[string]any {
    if s == nil {
        return nil
    }
    return s.AsMap()
}

// {{.Var}}ToStruct JSON值转proto结构
func {{.Var}}ToStruct(v *model.JsonValue) *structpb.Struct {
    if v == nil {
        return nil
    }
    m, ok := v.Data.(map[string]any)
    if !ok {
        return nil
    }
    s, err := structpb.NewStruct(m)
    if err != nil {
        return nil
    }
    return s
}

// toProto{{.Name}} 模型转proto
func toProto{{.Name}}(m *model.{{.Name}}) *proto.{{.Name}} {
    if m == nil {
        return nil
    }
    return &proto.{{.Name}}{
{{.ProtoFields}}
    }
}
