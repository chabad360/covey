// +build !live

package ui

import "html/template"

var (
	templates = make(map[string]*template.Template)
)

// GetTemplate returns a template from the template map.
func GetTemplate(name string) *template.Template {
	templates[name] = templatesF[name]()
	return templates[name]
}
