package {{.Package}}

// {{.FacadeName}} {{.Desc}}
// 使用示例:
//   facade.{{.FacadeName}}().Method()
func {{.FacadeName}}() *{{.FacadeVar}}Facade {
	return &{{.FacadeVar}}Facade{}
}

type {{.FacadeVar}}Facade struct{}

// todo 添加具体方法
// func (f *{{.FacadeVar}}Facade) Method() {
//     todo 方法实现
// }
