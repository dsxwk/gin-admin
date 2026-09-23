package model{{ if .Imports }}

{{ .Imports }}{{ end }}

const {{ .TableConst }} = "{{ .Table }}"

// {{.Struct}} {{.StructComment}}
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
func (*{{ .Struct }}) Connection() string {
    return "{{ .Connection }}"
}
{{- end}}
