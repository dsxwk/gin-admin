package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Article 请求验证
type Article struct {
	base.BaseRequest
	ID         int64       `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Uid        int64       `json:"uid" form:"uid" validate:"int" label:"用户id"`
	Title      string      `json:"title" form:"title" validate:"required" label:"标题"`
	Content    string      `json:"content" form:"content" validate:"required" label:"内容"`
	CategoryId int64       `json:"categoryId" form:"categoryId" validate:"int" label:"分类id"`
	DataSource int64       `json:"dataSource" form:"dataSource" validate:"int" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64       `json:"isPublish" form:"isPublish" validate:"int" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        interface{} `json:"tag" form:"tag" validate:"required" label:"标签"`
	PageListValidate
}

// ArticleCreate 文章创建验证
type ArticleCreate struct {
	Uid        int64       `json:"uid" form:"uid" validate:"int" label:"用户id"`
	Title      string      `json:"title" form:"title" validate:"required" label:"标题"`
	Content    string      `json:"content" form:"content" validate:"required" label:"内容"`
	CategoryId int64       `json:"categoryId" form:"categoryId" validate:"int" label:"分类id"`
	DataSource int64       `json:"dataSource" form:"dataSource" validate:"int" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64       `json:"isPublish" form:"isPublish" validate:"int" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        interface{} `json:"tag" form:"tag" validate:"required" label:"标签"`
}

// ArticleUpdate 文章更新验证
type ArticleUpdate struct {
	ID         int64       `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Uid        int64       `json:"uid" form:"uid" validate:"int" label:"用户id"`
	Title      string      `json:"title" form:"title" validate:"required" label:"标题"`
	Content    string      `json:"content" form:"content" validate:"required" label:"内容"`
	CategoryId int64       `json:"categoryId" form:"categoryId" validate:"int" label:"分类id"`
	DataSource int64       `json:"dataSource" form:"dataSource" validate:"int" label:"数据来源 1=文章库 2=自建"`
	IsPublish  int64       `json:"isPublish" form:"isPublish" validate:"int" label:"是否发布 0=待发布 1=已发布 2=已下架"`
	Tag        interface{} `json:"tag" form:"tag" validate:"required" label:"标签"`
}

// Validate 请求验证
func (s Article) Validate(data Article, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Article) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List": []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{
			"Uid",
			"Title",
			"Content",
			"CategoryId",
			"DataSource",
			"IsPublish",
			"Tag",
		},
		"Update": []string{
			"ID",
			"Uid",
			"Title",
			"Content",
			"CategoryId",
			"DataSource",
			"IsPublish",
			"Tag",
		},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Article) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"minLen":                       "{field} 长度不能少于 {min} 个字符",
		"maxLen":                       "{field} 长度不能超过 {max} 个字符",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s Article) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Uid":                       "用户id",
		"Title":                     "标题",
		"Content":                   "内容",
		"CategoryId":                "分类id",
		"DataSource":                "数据来源 1=文章库 2=自建",
		"IsPublish":                 "是否发布 1=已发布 2=未发布",
		"Tag":                       "标签",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
