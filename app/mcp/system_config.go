package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

// SystemConfigQueryArgs 系统配置查询参数
type SystemConfigQueryArgs struct {
	Key string `json:"key"`
}

// SystemConfigQueryTool 系统配置查询工具
type SystemConfigQueryTool struct{}

// NewSystemConfigQueryTool 创建系统配置查询工具
func NewSystemConfigQueryTool() *mcp.ToolAdapter[SystemConfigQueryArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "query_system_config",
			Description: "查询系统配置项,如站点名称、LOGO等",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"key": {Type: "string", Description: "配置键名,如site_name,留空查全部"},
				},
			},
		},
		(&SystemConfigQueryTool{}).Handle,
	)
}

// Handle 处理系统配置查询工具
func (t *SystemConfigQueryTool) Handle(ctx context.Context, _ int64, args *SystemConfigQueryArgs) (any, error) {
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	var configs []model.SystemConfig
	query := db.Model(&model.SystemConfig{})
	if args.Key != "" {
		query = query.Where("`key` = ?", args.Key)
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
