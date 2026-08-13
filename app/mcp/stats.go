package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

func init() { mcp.Register(&OperatorStats{}) }

// OperatorStats 操作日志统计工具
type OperatorStats struct{}

func (t *OperatorStats) Name() string { return "query_operator_stats" }
func (t *OperatorStats) Description() string {
	return "查询操作日志统计,如今日PV/UV、请求方法分布、状态码统计等"
}

func (t *OperatorStats) InputSchema() mcp.InputSchema {
	return mcp.InputSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"date": {Type: "string", Description: "统计日期,格式2026-08-12,默认今天"},
		},
	}
}

func (t *OperatorStats) Call(ctx context.Context, args map[string]any) (any, error) {
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	date, _ := args["date"].(string)
	if date == "" {
		return nil, nil
	}

	// 统计PV
	var pv int64
	db.Model(&model.OperatorLog{}).Where("DATE(created_at) = ?", date).Count(&pv)

	// 统计UV(按IP去重)
	var uv int64
	db.Model(&model.OperatorLog{}).Where("DATE(created_at) = ?", date).Distinct("ip").Count(&uv)

	// 统计请求方法分布
	type MethodCount struct {
		Method string `json:"method"`
		Count  int64  `json:"count"`
	}
	var methods []MethodCount
	db.Model(&model.OperatorLog{}).Select("method, COUNT(*) as count").
		Where("DATE(created_at) = ?", date).Group("method").Order("count DESC").Scan(&methods)

	// 统计状态码分布
	type StatusCount struct {
		StatusCode int64 `json:"statusCode"`
		Count      int64 `json:"count"`
	}
	var statuses []StatusCount
	db.Model(&model.OperatorLog{}).Select("status_code, COUNT(*) as count").
		Where("DATE(created_at) = ?", date).Group("status_code").Order("count DESC").Scan(&statuses)

	return map[string]any{
		"date":     date,
		"pv":       pv,
		"uv":       uv,
		"methods":  methods,
		"statuses": statuses,
	}, nil
}
