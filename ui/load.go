// +build !live

package ui

import "html/template"

var (
	templates map[string]*template.Template
)

func getTemplate(name string) *template.Template {
	return templates[name]
}

func init() {
	templates = map[string]*template.Template{
		"base":  templatesF["base"](),
		"login": templatesF["login"](),
	}
}
