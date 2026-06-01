package {{.Package}}

import (
    "gin/common/base"
)

type {{.Name}}Service struct {
    base.BaseService
}

// {{.Function}} {{.Description}}
func (s *{{.Name}}Service) {{.Function}}() {
    // Define your function here
}