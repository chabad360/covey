// +build !live

package ui

import "html/template"

var (
	templates map[string]*template.Template
)

// GetTemplate returns a template from the template map.
func GetTemplate(name string) *template.Template {
	return templates[name]
}

func init() {
	templates = map[string]*template.Template{
		"dashboard":   templatesF["dashboard"](),
		"tasksAll":    templatesF["tasksAll"](),
		"tasksSingle": templatesF["tasksSingle"](),
		"login":       templatesF["login"](),
	}
}
