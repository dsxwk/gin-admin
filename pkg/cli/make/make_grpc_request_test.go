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

// TestGrpcRequestTemplateUsesFields 验证gRPC请求模板生成Fields校验
func TestGrpcRequestTemplateUsesFields(t *testing.T) {
	maker := &MakeGrpcRequest{}
	fields := []GrpcRequestField{
		{Name: "Name", JSON: "name", Type: "string", Label: "名称", Required: true},
		{Name: "Age", JSON: "age", Type: "int32", Label: "年龄"},
	}
	maker.formatGrpcRequestFields(fields)
	formattedId, formattedUpdateId := maker.formatGrpcIdFields(fields)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate gRPC request template test")
	}
	templatePath := filepath.Join(filepath.Dir(currentFile), "..", "template", "grpc_request.tpl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		t.Fatalf("parse gRPC request template: %v", err)
	}

	data := struct {
		Package           string
		Name              string
		Var               string
		Description       string
		Fields            []GrpcRequestField
		FormattedId       string
		FormattedUpdateId string
		Translates        []string
	}{
		Package:           "request",
		Name:              "Demo",
		Var:               "demo",
		Description:       "示例",
		Fields:            fields,
		FormattedId:       formattedId,
		FormattedUpdateId: formattedUpdateId,
		Translates: []string{
			`"Id": "ID",`,
			`"Name": "名称",`,
		},
	}

	var output strings.Builder
	if err = tmpl.Execute(&output, data); err != nil {
		t.Fatalf("execute gRPC request template: %v", err)
	}
	content := output.String()
	if _, err = parser.ParseFile(token.NewFileSet(), "user.go", content, parser.AllErrors); err != nil {
		t.Fatalf("generated gRPC request is invalid Go: %v\n%s", err, content)
	}

	for _, expected := range []string{
		"func (s DemoCreateRequest) Validate(fields validatex.Fields) validatex.Errors",
		`fields.Str("Name", s.Name).Required().Errors(),`,
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("generated gRPC request missing %q:\n%s", expected, content)
		}
	}
	for _, unexpected := range []string{
		"github.com/gookit/validate",
		"ConfigValidation",
		" validate:",
	} {
		if strings.Contains(content, unexpected) {
			t.Fatalf("generated gRPC request contains %q:\n%s", unexpected, content)
		}
	}
}
