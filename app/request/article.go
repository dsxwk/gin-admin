package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Article 请求验证
type Article struct {
	base.Request
	ID         int64  `json:"id" form:"id" label:"ID"`
	Title      string `json:"title" form:"title" label:"标题"`
	Content    string `json:"content" form:"content" label:"内容"`
	CategoryId int64  `json:"categoryId" form:"categoryId" label:"分类id"`
	DataSource int64  `json:"dataSource" form:"dataSource" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64  `json:"isPublish" form:"isPublish" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        any    `json:"tag" form:"tag" label:"标签"`
	PageListValidate
}

// ArticleCreate 文章创建验证
type ArticleCreate struct {
	Title      string `json:"title" form:"title" label:"标题"`
	Content    string `json:"content" form:"content" label:"内容"`
	CategoryId int64  `json:"categoryId" form:"categoryId" label:"分类id"`
	DataSource int64  `json:"dataSource" form:"dataSource" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64  `json:"isPublish" form:"isPublish" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        any    `json:"tag" form:"tag" label:"标签"`
}

// ArticleUpdate 文章更新验证
type ArticleUpdate struct {
	Title      string `json:"title" form:"title" label:"标题"`
	Content    string `json:"content" form:"content" label:"内容"`
	CategoryId int64  `json:"categoryId" form:"categoryId" label:"分类id"`
	DataSource int64  `json:"dataSource" form:"dataSource" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64  `json:"isPublish" form:"isPublish" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        any    `json:"tag" form:"tag" label:"标签"`
}

// ValidateWithScene 请求验证
func (s Article) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Str("Title", s.Title).Required().Errors(),
			fields.Str("Content", s.Content).Required().Errors(),
			requiredValue(fields, "Tag", s.Tag),
		)
	case "Update":
		return validatex.Check(
			fields.Int("ID", s.ID).Required().Gt(0).Errors(),
			fields.Str("Title", s.Title).Required().Errors(),
			fields.Str("Content", s.Content).Required().Errors(),
			requiredValue(fields, "Tag", s.Tag),
		)
	case "Detail", "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Article) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Article) Translates() validatex.Translates {
	return validatex.Translates{
		"ID":                        "ID",
		"Title":                     "标题",
		"Content":                   "内容",
		"CategoryId":                "分类id",
		"DataSource":                "数据来源",
		"IsPublish":                 "是否发布",
		"Tag":                       "标签",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
