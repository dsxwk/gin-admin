package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

// OperatorStatsArgs 操作日志统计参数
type OperatorStatsArgs struct {
	Date string `json:"date"`
}

// OperatorStatsTool 操作日志统计工具
type OperatorStatsTool struct{}

// NewOperatorStatsTool 创建操作日志统计工具
func NewOperatorStatsTool() *mcp.ToolAdapter[OperatorStatsArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "query_operator_stats",
			Description: "查询操作日志统计,如今日PV/UV、请求方法分布、状态码统计等",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"date": {Type: "string", Description: "统计日期,格式2026-08-12,默认今天"},
				},
			},
		},
		(&OperatorStatsTool{}).Handle,
	)
}

// Handle 处理操作日志统计工具
func (t *OperatorStatsTool) Handle(ctx context.Context, _ int64, args *OperatorStatsArgs) (any, error) {
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}
	if args.Date == "" {
		return nil, nil
	}

	// 统计PV
	var pv int64
	db.Model(&model.OperatorLog{}).Where("DATE(created_at) = ?", args.Date).Count(&pv)

	// 统计UV(按IP去重)
	var uv int64
	db.Model(&model.OperatorLog{}).Where("DATE(created_at) = ?", args.Date).Distinct("ip").Count(&uv)

	// 统计请求方法分布
	type MethodCount struct {
		Method string `json:"method"`
		Count  int64  `json:"count"`
	}
	var methods []MethodCount
	db.Model(&model.OperatorLog{}).Select("method, COUNT(*) as count").
		Where("DATE(created_at) = ?", args.Date).Group("method").Order("count DESC").Scan(&methods)

	// 统计状态码分布
	type StatusCount struct {
		StatusCode int64 `json:"statusCode"`
		Count      int64 `json:"count"`
	}
	var statuses []StatusCount
	db.Model(&model.OperatorLog{}).Select("status_code, COUNT(*) as count").
		Where("DATE(created_at) = ?", args.Date).Group("status_code").Order("count DESC").Scan(&statuses)

	return map[string]any{
		"date":     args.Date,
		"pv":       pv,
		"uv":       uv,
		"methods":  methods,
		"statuses": statuses,
	}, nil
}
