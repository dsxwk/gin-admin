package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

func init() { mcp.Register(&DictQuery{}) }

// DictQuery 字典查询工具
type DictQuery struct{}

func (t *DictQuery) Name() string { return "query_dict" }
func (t *DictQuery) Description() string {
	return "根据字典类型查询字典项,如性别、状态等枚举值"
}

func (t *DictQuery) InputSchema() mcp.InputSchema {
	return mcp.InputSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"type":  {Type: "string", Description: "字典类型编码,如gender、status"},
			"label": {Type: "string", Description: "按字典标签模糊搜索"},
		},
		Required: []string{"type"},
	}
}

func (t *DictQuery) Call(ctx context.Context, args map[string]any) (any, error) {
	typeName, _ := args["type"].(string)
	if typeName == "" {
		return nil, nil
	}

	db := facade.DB().WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	var dicts []model.Dict
	query := db.Model(&model.Dict{}).Where("name = ?", typeName).Order("sort ASC, id ASC")

	if label, ok := args["label"].(string); ok && label != "" {
		query = query.Where("title LIKE ?", "%"+label+"%")
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
	return map[string]any{"type": typeName, "total": len(items), "list": items}, nil
}
