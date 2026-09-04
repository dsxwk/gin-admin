package {{.Package}}

import (
    "gin/common/errcode"

    "github.com/gookit/validate"
)

// {{.Name}}Request {{.Description}}请求
type {{.Name}}Request struct {
    {{.FormattedId}}
}

// {{.Name}}ListRequest {{.Description}}列表请求
type {{.Name}}ListRequest struct {
    Page     int            `json:"page" validate:"required|int|gt:0" label:"页码"`
    PageSize int            `json:"pageSize" validate:"required|int|gt:0" label:"每页数量"`
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
    Ids []int32 `json:"ids" validate:"required|minLen:1" label:"ID列表"`
}

// Validate {{.Description}}请求验证
func (s {{.Name}}Request) Validate(data {{.Name}}Request, scene string) error {
    v := validate.Struct(data, scene)
    if !v.Validate(scene) {
        return errcode.ArgsError().WithMsg(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
func (s {{.Name}}Request) ConfigValidation(v *validate.Validation) {
    v.WithScenes(validate.SValues{
        "Detail": []string{"Id"},
        "Delete": []string{"Id"},
    })
}

// Messages 验证器错误消息
func (s {{.Name}}Request) Messages() map[string]string {
    return {{.Var}}Messages()
}

// Translates 字段翻译
func (s {{.Name}}Request) Translates() map[string]string {
    return {{.Var}}Translates()
}

// Validate 列表请求验证
func (s {{.Name}}ListRequest) Validate() error {
    v := validate.Struct(s, "List")
    if !v.Validate("List") {
        return errcode.ArgsError().WithMsg(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
func (s {{.Name}}ListRequest) ConfigValidation(v *validate.Validation) {
    v.WithScenes(validate.SValues{
        "List": []string{"Page", "PageSize"},
    })
}

// Messages 验证器错误消息
func (s {{.Name}}ListRequest) Messages() map[string]string {
    return {{.Var}}Messages()
}

// Translates 字段翻译
func (s {{.Name}}ListRequest) Translates() map[string]string {
    return {{.Var}}Translates()
}

// Validate 创建请求验证
func (s {{.Name}}CreateRequest) Validate() error {
    v := validate.Struct(s, "Create")
    if !v.Validate("Create") {
        return errcode.ArgsError().WithMsg(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
func (s {{.Name}}CreateRequest) ConfigValidation(v *validate.Validation) {
    v.WithScenes(validate.SValues{
        "Create": []string{
{{- range .CreateScene}}
            "{{.}}",
{{- end}}
        },
    })
}

// Messages 验证器错误消息
func (s {{.Name}}CreateRequest) Messages() map[string]string {
    return {{.Var}}Messages()
}

// Translates 字段翻译
func (s {{.Name}}CreateRequest) Translates() map[string]string {
    return {{.Var}}Translates()
}

// Validate 更新请求验证
func (s {{.Name}}UpdateRequest) Validate() error {
    v := validate.Struct(s, "Update")
    if !v.Validate("Update") {
        return errcode.ArgsError().WithMsg(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
func (s {{.Name}}UpdateRequest) ConfigValidation(v *validate.Validation) {
    v.WithScenes(validate.SValues{
        "Update": []string{
            "Id",
{{- range .Fields}}
            "{{.Name}}",
{{- end}}
        },
    })
}

// Messages 验证器错误消息
func (s {{.Name}}UpdateRequest) Messages() map[string]string {
    return {{.Var}}Messages()
}

// Translates 字段翻译
func (s {{.Name}}UpdateRequest) Translates() map[string]string {
    return {{.Var}}Translates()
}

// Validate 批量删除请求验证
func (s {{.Name}}BatchDeleteRequest) Validate() error {
    v := validate.Struct(s, "BatchDelete")
    if !v.Validate("BatchDelete") {
        return errcode.ArgsError().WithMsg(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
func (s {{.Name}}BatchDeleteRequest) ConfigValidation(v *validate.Validation) {
    v.WithScenes(validate.SValues{
        "BatchDelete": []string{"Ids"},
    })
}

// Messages 验证器错误消息
func (s {{.Name}}BatchDeleteRequest) Messages() map[string]string {
    return {{.Var}}Messages()
}

// Translates 字段翻译
func (s {{.Name}}BatchDeleteRequest) Translates() map[string]string {
    return {{.Var}}Translates()
}

// {{.Var}}Messages {{.Description}}公共错误消息
func {{.Var}}Messages() map[string]string {
    return validate.MS{
        "required": "字段 {field} 必填",
        "int":      "字段 {field} 必须为整数",
        "float":    "字段 {field} 必须为数字",
        "gt":       "字段 {field} 需大于 0",
        "minLen":   "字段 {field} 不能为空",
    }
}

// {{.Var}}Translates {{.Description}}公共字段翻译
func {{.Var}}Translates() map[string]string {
    return validate.MS{
{{- range .Translates}}
        {{.}}
{{- end}}
    }
}
