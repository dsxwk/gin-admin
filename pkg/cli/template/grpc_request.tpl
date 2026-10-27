package {{.Package}}

import "github.com/dsxwk/validatex"

// {{.Name}}Request {{.Description}}请求
type {{.Name}}Request struct {
	{{.FormattedId}}
}

// {{.Name}}ListRequest {{.Description}}列表请求
type {{.Name}}ListRequest struct {
	Page     int            `json:"page" label:"页码"`
	PageSize int            `json:"pageSize" label:"每页数量"`
	NotPage  bool           `json:"notPage" label:"不分页"`
	Search   map[string]any `json:"search" label:"搜索条件"`
	Sort     map[string]any `json:"sort" label:"排序条件"`
}

// {{.Name}}CreateRequest {{.Description}}创建请求
type {{.Name}}CreateRequest struct {
{{- range .Fields}}
	{{.FormattedField}}
{{- end}}
}

// {{.Name}}UpdateRequest {{.Description}}更新请求
type {{.Name}}UpdateRequest struct {
	{{.FormattedUpdateId}}
{{- range .Fields}}
	{{.FormattedUpdateField}}
{{- end}}
}

// {{.Name}}BatchDeleteRequest {{.Description}}批量删除请求
type {{.Name}}BatchDeleteRequest struct {
	Ids []int32 `json:"ids" label:"ID列表"`
}

// ValidateWithScene {{.Description}}请求验证
func (s {{.Name}}Request) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Detail", "Delete":
		return fields.Int("Id", s.Id).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Validate 列表请求验证
func (s {{.Name}}ListRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.When(!s.NotPage, func() validatex.Errors {
		return validatex.Check(
			fields.Int("Page", s.Page).Required().Gt(0).Errors(),
			fields.Int("PageSize", s.PageSize).Required().Gt(0).Errors(),
		)
	})
}

// Validate 创建请求验证
func (s {{.Name}}CreateRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
{{- range .Fields}}
{{- if .Validation}}
		{{.Validation}}
{{- end}}
{{- end}}
	)
}

// Validate 更新请求验证
func (s {{.Name}}UpdateRequest) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Int("Id", s.Id).Required().Gt(0).Errors()
}

// Validate 批量删除请求验证
func (s {{.Name}}BatchDeleteRequest) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Slice("Ids", s.Ids).Required().Each("gt", "", func(id int32) bool {
		return id > 0
	}, int32(0)).Errors()
}

// Messages {{.Description}}错误消息
func (s {{.Name}}Request) Messages() validatex.Messages {
	return {{.Var}}Messages()
}

// Translates {{.Description}}字段翻译
func (s {{.Name}}Request) Translates() validatex.Translates {
	return {{.Var}}Translates()
}

// Messages {{.Description}}列表错误消息
func (s {{.Name}}ListRequest) Messages() validatex.Messages {
	return {{.Var}}Messages()
}

// Translates {{.Description}}列表字段翻译
func (s {{.Name}}ListRequest) Translates() validatex.Translates {
	return {{.Var}}Translates()
}

// Messages {{.Description}}创建错误消息
func (s {{.Name}}CreateRequest) Messages() validatex.Messages {
	return {{.Var}}Messages()
}

// Translates {{.Description}}创建字段翻译
func (s {{.Name}}CreateRequest) Translates() validatex.Translates {
	return {{.Var}}Translates()
}

// Messages {{.Description}}更新错误消息
func (s {{.Name}}UpdateRequest) Messages() validatex.Messages {
	return {{.Var}}Messages()
}

// Translates {{.Description}}更新字段翻译
func (s {{.Name}}UpdateRequest) Translates() validatex.Translates {
	return {{.Var}}Translates()
}

// Messages {{.Description}}批量删除错误消息
func (s {{.Name}}BatchDeleteRequest) Messages() validatex.Messages {
	return {{.Var}}Messages()
}

// Translates {{.Description}}批量删除字段翻译
func (s {{.Name}}BatchDeleteRequest) Translates() validatex.Translates {
	return {{.Var}}Translates()
}

// {{.Var}}Messages {{.Description}}公共错误消息
func {{.Var}}Messages() validatex.Messages {
	return validatex.Messages{
		"required": "字段 {field} 必填",
		"gt":       "字段 {field} 需大于 0",
	}
}

// {{.Var}}Translates {{.Description}}公共字段翻译
func {{.Var}}Translates() validatex.Translates {
	return validatex.Translates{
{{- range .Translates}}
		{{.}}
{{- end}}
	}
}
