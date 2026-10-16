package mcp

import (
	"context"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/serviceprovider/mcp"
)

// UserSearchArgs 用户搜索参数
type UserSearchArgs struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

// NewUserSearchTool 创建用户搜索工具
func NewUserSearchTool() *mcp.ToolAdapter[UserSearchArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "search_users",
			Description: "根据姓名或用户名搜索系统用户",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"keyword": {Type: "string", Description: "搜索关键词,匹配姓名或用户名"},
					"limit":   {Type: "integer", Description: "返回条数,默认20,最大100"},
				},
				Required: []string{"keyword"},
			},
		},
		handleUserSearch,
	)
}

// handleUserSearch 处理用户搜索工具
func handleUserSearch(ctx context.Context, _ int64, args *UserSearchArgs) (any, error) {
	if args.Keyword == "" {
		return nil, nil
	}

	limit := 20
	if args.Limit > 0 {
		limit = min(args.Limit, 100)
	}

	var users []model.User
	like := "%" + args.Keyword + "%"
	db := facade.DB("mysql").WithContext(ctx)
	if db == nil {
		return nil, nil
	}

	err := db.Model(&model.User{}).
		Where("(full_name LIKE ? OR username LIKE ?)", like, like).
		Preload("MainDept", "is_main = ?", 1).
		Preload("MainDept.Dept").
		Limit(limit).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}

	type Item struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Status   int64  `json:"status"`
		Dept     string `json:"dept"`
	}
	items := make([]Item, len(users))
	for i, u := range users {
		dept := ""
		if u.MainDept != nil && u.MainDept.Dept != nil {
			dept = u.MainDept.Dept.Name
		}
		items[i] = Item{ID: u.ID, Username: u.Username, FullName: u.FullName, Email: u.Email, Status: u.Status, Dept: dept}
	}
	return map[string]any{"keyword": args.Keyword, "total": len(items), "list": items}, nil
}
