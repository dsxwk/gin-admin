syntax = "proto3";

package {{.Package}};

import "google/protobuf/struct.proto";
import "grpc/proto/base.proto";

option go_package = "{{.GoPackage}}";

// 请求消息
message {{.Name}}Request {
  int32 id = 1;
}

message {{.Name}}ListRequest {
  int32 page = 1;
  int32 pageSize = 2;
  bool notPage = 3;
  google.protobuf.Struct search = 4;
  google.protobuf.Struct sort = 5;
}

message {{.Name}}CreateRequest {
{{.CreateFields}}
}

message {{.Name}}UpdateRequest {
  int32 id = 1;
  google.protobuf.Struct data = 2;
}

message {{.Name}} {
  int32 id = 1;
{{.UserFields}}
}

message {{.Name}}ListResponse {
  int32 total = 1;
  int32 page = 2;
  int32 pageSize = 3;
  repeated {{.Name}} list = 4;
}

// 服务定义
service {{.ServiceName}} {
  rpc Detail({{.Name}}Request) returns ({{.Name}});
  rpc List({{.Name}}ListRequest) returns ({{.Name}}ListResponse);
  rpc Create({{.Name}}CreateRequest) returns ({{.Name}});
  rpc Update({{.Name}}UpdateRequest) returns ({{.Name}});
  rpc Delete({{.Name}}Request) returns (EmptyResponse);
}
