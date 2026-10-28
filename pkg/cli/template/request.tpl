package {{.Package}}

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// {{.StructName}} {{.Description}}
type {{.StructName}} struct {
	base.Request
{{- range .Fields}}
{{.FormattedField}}
{{- end}}
	PageListValidate
}

// ValidateWithScene 请求验证
func (s {{.StructName}}) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return validatex.When(!s.NotPage, func() validatex.Errors {
			return validatex.Check(
				fields.Int("page", s.Page).Required().Gt(0).Errors(),
				fields.Int("pageSize", s.PageSize).Required().Gt(0).Errors(),
			)
		})
	case "Create":
		return validatex.Check(
{{- range .Fields}}
{{- if and .Validation (not .IsID)}}
			{{.Validation}}
{{- end}}
{{- end}}
		)
	case "Update":
		return validatex.Check(
			fields.Int("{{.IDJSON}}", s.ID).Required().Gt(0).Errors(),
{{- range .Fields}}
{{- if and .Validation (not .IsID)}}
			{{.Validation}}
{{- end}}
{{- end}}
		)
	case "Detail", "Delete":
		return fields.Int("{{.IDJSON}}", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s {{.StructName}}) Messages() validatex.Messages {
	return validatex.Messages{
		"required": "字段 {field} 必填",
		"gt":       "字段 {field} 需大于 0",
		"min":      "{field} 长度不能少于 {min} 个字符",
		"max":      "{field} 长度不能超过 {max} 个字符",
	}
}

// Translates 字段翻译
func (s {{.StructName}}) Translates() validatex.Translates {
	return validatex.Translates{
{{- range .FormattedTranslates}}
		{{.}}
{{- end}}
	}
}
