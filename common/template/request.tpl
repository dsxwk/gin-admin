package {{.Package}}

import (
    "errors"
    "gin/common/base"
    "github.com/gookit/validate"
)

// {{.StructName}} {{.Description}}
type {{.StructName}} struct {
    base.BaseRequest
{{- range .Fields}}
{{.FormattedField}}
{{- end}}
    PageListValidate
}

// Validate 请求验证
func (s {{.StructName}}) Validate(data {{.StructName}}, scene string) error {
    v := validate.Struct(data, scene)
    if !v.Validate(scene) {
        return errors.New(v.Errors.One())
    }
    return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s {{.StructName}}) ConfigValidation(v *validate.Validation) {
    scenes := validate.SValues{
        "list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
        "create": []string{ {{.CreateScene}} },
        "update": []string{ {{.UpdateScene}} },
        "detail": []string{"ID"},
        "delete": []string{"ID"},
    }
    v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s {{.StructName}}) Messages() map[string]string {
    return validate.MS{
        "required":                     "字段 {field} 必填",
        "int":                          "字段 {field} 必须为整数",
        "gt":                           "字段 {field} 必须大于 0",
        "minLen":                       "{field} 长度不能少于 {min} 个字符",
        "maxLen":                       "字段 {field} 长度不能超过 255",
        "PageListValidate.Page.gt":     "页码必须大于 0",
        "PageListValidate.PageSize.gt": "每页数量必须大于 0",
    }
}

// Translates 字段翻译
func (s {{.StructName}}) Translates() map[string]string {
    return validate.MS{
{{- range .FormattedTranslates}}
        {{.}}
{{- end}}
    }
}
