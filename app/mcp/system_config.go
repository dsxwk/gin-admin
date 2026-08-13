package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

func init() { mcp.Register(&SystemConfigQuery{}) }

// SystemConfigQuery 系统配置查询工具
type SystemConfigQuery struct{}

func (t *SystemConfigQuery) Name() string { return "query_system_config" }
func (t *SystemConfigQuery) Description() string {
	return "查询系统配置项,如站点名称、LOGO等"
}

func (t *SystemConfigQuery) InputSchema() mcp.InputSchema {
	return mcp.InputSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"key": {Type: "string", Description: "配置键名,如site_name,留空查全部"},
		},
	}
}

func (t *SystemConfigQuery) Call(ctx context.Context, args map[string]any) (any, error) {
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	var configs []model.SystemConfig
	query := db.Model(&model.SystemConfig{})

	if key, ok := args["key"].(string); ok && key != "" {
		query = query.Where("`key` = ?", key)
	}

	if err := query.Find(&configs).Error; err != nil {
		return nil, err
	}

	type Item struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Value string `json:"value"`
		Type  int64  `json:"type"`
	}
	items := make([]Item, len(configs))
	for i, c := range configs {
		items[i] = Item{Key: c.Key, Name: c.Name, Value: c.DefaultValue, Type: c.Type}
	}
	return map[string]any{"total": len(items), "list": items}, nil
}
