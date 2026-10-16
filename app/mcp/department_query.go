package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
	"strconv"
)

// DepartmentQueryArgs 部门查询参数
type DepartmentQueryArgs struct {
	Pid       *int   `json:"pid"`
	Name      string `json:"name"`
	WithUsers bool   `json:"with_users"`
}

// DepartmentQueryTool 部门查询工具
type DepartmentQueryTool struct{}

// NewDepartmentQueryTool 创建部门查询工具
func NewDepartmentQueryTool() *mcp.ToolAdapter[DepartmentQueryArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "query_departments",
			Description: "查询部门列表或指定上级的子部门,返回部门树结构",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"pid":        {Type: "integer", Description: "上级部门ID,0查全部顶级部门"},
					"name":       {Type: "string", Description: "按部门名称模糊搜索"},
					"with_users": {Type: "boolean", Description: "是否包含部门成员数量,默认false"},
				},
			},
		},
		(&DepartmentQueryTool{}).Handle,
	)
}

// Handle 处理部门查询工具
func (t *DepartmentQueryTool) Handle(ctx context.Context, _ int64, args *DepartmentQueryArgs) (any, error) {
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	query := db.Model(&model.Department{}).Order("sort ASC, id ASC")
	if args.Pid != nil {
		query = query.Where("pid = ?", *args.Pid)
	}
	if args.Name != "" {
		query = query.Where("name LIKE ?", "%"+args.Name+"%")
	}

	var depts []model.Department
	if err := query.Find(&depts).Error; err != nil {
		return nil, err
	}

	type DeptItem struct {
		ID     int64  `json:"id"`
		Pid    int64  `json:"pid"`
		Name   string `json:"name"`
		Status int64  `json:"status"`
		Sort   int64  `json:"sort"`
		Path   string `json:"path"`
	}
	items := make([]DeptItem, len(depts))
	for i, d := range depts {
		items[i] = DeptItem{ID: d.ID, Pid: d.Pid, Name: d.Name, Status: d.Status, Sort: d.Sort, Path: strconv.FormatInt(d.ID, 10)}
	}
	return map[string]any{"total": len(items), "list": items}, nil
}
