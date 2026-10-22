package registry

import (
	"gin/pkg/cli/base"
	"gin/pkg/cli/db"
	"gin/pkg/cli/event"
	"gin/pkg/cli/grpc"
	"gin/pkg/cli/job"
	"gin/pkg/cli/jwt"
	mk "gin/pkg/cli/make"
	"gin/pkg/cli/mcp"
	"gin/pkg/cli/migrate"
	"gin/pkg/cli/permission"
	"gin/pkg/cli/queue"
	"gin/pkg/cli/route"
)

// All 获取全部内置CLI命令
func All() []base.Handler {
	return []base.Handler{
		// 数据库
		db.NewInit(),
		db.NewSeed(),
		// 迁移
		migrate.NewFresh(),
		migrate.NewMigrate(),
		migrate.NewRollback(),
		migrate.NewRefresh(),
		migrate.NewStatus(),
		// 事件
		event.NewEventList(),
		event.NewEventListenerList(),
		// gRPC
		grpc.NewGrpcGen(),
		// 任务调度
		job.NewJobList(),
		job.NewJobClear(),
		// JWT
		jwt.NewJwt(),
		// 命令行生成
		mk.NewMakeConsole(),
		mk.NewMakeController(),
		mk.NewMakeDocs(),
		mk.NewMakeEnum(),
		mk.NewMakeErrCode(),
		mk.NewMakeEs(),
		mk.NewMakeEvent(),
		mk.NewMakeFacade(),
		mk.NewMakeGrpcModel(),
		mk.NewMakeGrpcProto(),
		mk.NewMakeGrpcRequest(),
		mk.NewMakeGrpcService(),
		mk.NewMakeJob(),
		mk.NewMakeListener(),
		mk.NewMakeMcp(),
		mk.NewMakeMiddleware(),
		mk.NewMakeMigration(),
		mk.NewMakeModel(),
		mk.NewMakeModelOld(),
		mk.NewMakeProvider(),
		mk.NewMakeQueue(),
		mk.NewMakeRequest(),
		mk.NewMakeRouter(),
		mk.NewMakeSeeder(),
		mk.NewMakeService(),
		mk.NewMakeWs(),
		// MCP
		mcp.NewMcpList(),
		// 权限
		permission.NewPermissionSync(),
		// 队列
		queue.NewConsumerList(),
		queue.NewProducerList(),
		// 路由
		route.NewRouteList(),
	}
}
