package {{.Package}}

type {{.Struct}}Event struct {
	// define event fields here
}

func (e {{.Struct}}Event) Name() string {
	return "{{.Name}}"
}

func (e {{.Struct}}Event) Description() string {
	return "{{.Description}}"
}
