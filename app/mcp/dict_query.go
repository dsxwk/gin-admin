package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

// DictQueryArgs 字典查询参数
type DictQueryArgs struct {
	Type  string `json:"type"`
	Label string `json:"label"`
}

// NewDictQueryTool 创建字典查询工具
func NewDictQueryTool() *mcp.ToolAdapter[DictQueryArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "query_dict",
			Description: "根据字典类型查询字典项,如性别、状态等枚举值",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"type":  {Type: "string", Description: "字典类型编码,如gender、status"},
					"label": {Type: "string", Description: "按字典标签模糊搜索"},
				},
				Required: []string{"type"},
			},
		},
		handleDictQuery,
	)
}

// handleDictQuery 处理字典查询工具
func handleDictQuery(ctx context.Context, _ int64, args *DictQueryArgs) (any, error) {
	if args.Type == "" {
		return nil, nil
	}

	db := facade.DB().WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	var dicts []model.Dict
	query := db.Model(&model.Dict{}).Where("name = ?", args.Type).Order("sort ASC, id ASC")
	if args.Label != "" {
		query = query.Where("title LIKE ?", "%"+args.Label+"%")
	}

	if err := query.Find(&dicts).Error; err != nil {
		return nil, err
	}

	type Item struct {
		ID    int64  `json:"id"`
		Label string `json:"label"`
		Value string `json:"value"`
		Sort  int64  `json:"sort"`
	}
	items := make([]Item, len(dicts))
	for i, d := range dicts {
		items[i] = Item{ID: d.ID, Label: d.Title, Value: d.Value, Sort: d.Sort}
	}
	return map[string]any{"type": args.Type, "total": len(items), "list": items}, nil
}
