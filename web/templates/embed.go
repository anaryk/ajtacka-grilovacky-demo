package templates

import (
	"embed"
	"html/template"
)

// Embed all HTML templates in the "web/templates" directory
//
//go:embed *.html
var templatesFS embed.FS

// NewTemplate function creates a new template from the embedded files
func NewTemplate(filenames ...string) (*template.Template, error) {
	tmpl := template.New("")
	return tmpl.ParseFS(templatesFS, filenames...)
}
