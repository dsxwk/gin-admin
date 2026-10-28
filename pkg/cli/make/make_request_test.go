package make

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"
)

// TestRequestTemplateUsesFields 验证请求模板生成Fields场景校验
func TestRequestTemplateUsesFields(t *testing.T) {
	maker := &MakeRequest{}
	fields := []Field{
		{Name: "ID", JSON: "id", Type: "int64", Label: "ID", Required: true, IsID: true},
		{Name: "Name", JSON: "name", Type: "string", Label: "名称", Required: true},
		{Name: "Age", JSON: "age", Type: "int64", Label: "年龄"},
	}
	maker.formatFields(fields)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate request template test")
	}
	templatePath := filepath.Join(filepath.Dir(currentFile), "..", "template", "request.tpl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		t.Fatalf("parse request template: %v", err)
	}

	data := struct {
		Package             string
		StructName          string
		Description         string
		Fields              []Field
		FormattedTranslates []string
		IDJSON              string
	}{
		Package:     "request",
		StructName:  "Demo",
		Description: "示例请求",
		Fields:      fields,
		IDJSON:      "id",
		FormattedTranslates: []string{
			`"name": "名称",`,
			`"page": "页码",`,
			`"pageSize": "每页数量",`,
		},
	}

	var output strings.Builder
	if err = tmpl.Execute(&output, data); err != nil {
		t.Fatalf("execute request template: %v", err)
	}
	content := output.String()
	if _, err = parser.ParseFile(token.NewFileSet(), "request.go", content, parser.AllErrors); err != nil {
		t.Fatalf("generated request is invalid Go: %v\n%s", err, content)
	}

	for _, expected := range []string{
		"func (s Demo) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors",
		`fields.Str("name", s.Name).Required().Errors(),`,
		`fields.Int("age", s.Age).Errors(),`,
		`fields.Int("id", s.ID).Required().Gt(0).Errors(),`,
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("generated request missing %q:\n%s", expected, content)
		}
	}
	for _, unexpected := range []string{
		"github.com/gookit/validate",
		"ConfigValidation",
		" validate:",
	} {
		if strings.Contains(content, unexpected) {
			t.Fatalf("generated request contains %q:\n%s", unexpected, content)
		}
	}
}
