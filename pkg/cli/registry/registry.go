package registry

import (
	"gin/common/base"
	"gin/pkg/cli/db"
	"gin/pkg/cli/event"
	"gin/pkg/cli/grpc"
	"gin/pkg/cli/job"
	"gin/pkg/cli/jwt"
	mk "gin/pkg/cli/make"
	"gin/pkg/cli/mcp"
	"gin/pkg/cli/permission"
	"gin/pkg/cli/queue"
	"gin/pkg/cli/route"
)

// All 获取全部内置CLI命令
func All() []base.Command {
	return []base.Command{
		// 数据库
		&db.Migrate{},
		&db.Rollback{},
		&db.Seed{},
		// 事件
		&event.EventList{},
		&event.EventListenerList{},
		// gRPC
		&grpc.GrpcGen{},
		// 任务调度
		&job.JobList{},
		&job.JobClear{},
		// JWT
		&jwt.Jwt{},
		// 命令行生成
		&mk.MakeConsole{},
		&mk.MakeController{},
		&mk.MakeDocs{},
		&mk.MakeEnum{},
		&mk.MakeErrCode{},
		&mk.MakeEs{},
		&mk.MakeEvent{},
		&mk.MakeFacade{},
		&mk.MakeGrpcModel{},
		&mk.MakeGrpcProto{},
		&mk.MakeGrpcRequest{},
		&mk.MakeGrpcService{},
		&mk.MakeJob{},
		&mk.MakeListener{},
		&mk.MakeMcp{},
		&mk.MakeMiddleware{},
		&mk.MakeMigration{},
		&mk.MakeModel{},
		&mk.MakeModelOld{},
		&mk.MakeProvider{},
		&mk.MakeQueue{},
		&mk.MakeRequest{},
		&mk.MakeRouter{},
		&mk.MakeSeed{},
		&mk.MakeService{},
		&mk.MakeWs{},
		// MCP
		&mcp.McpList{},
		// 权限
		&permission.PermissionSync{},
		// 队列
		&queue.ConsumerList{},
		&queue.ProducerList{},
		// 路由
		&route.RouteList{},
	}
}
