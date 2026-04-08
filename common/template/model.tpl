package model{{ if .Imports }}

{{ .Imports }}{{ end }}

const {{ .TableConst }} = "{{ .Table }}"

type {{ .Struct }} struct {
{{- range .Fields }}
	{{ . }}
{{- end }}
}

func (*{{ .Struct }}) TableName() string {
	return {{ .TableConst }}
}
{{- if .Connection}}

// Connection 数据库连接名称
func (m *{{ .Struct }}) Connection() string {
    return "{{ .Connection }}"
}
{{- end}}
