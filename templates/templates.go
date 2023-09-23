package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var templates embed.FS

func ParsedTemplateOrPanic(file ...string) *template.Template {
	return template.Must(template.ParseFS(templates, file...))
}
