// +build live

package ui

import "html/template"

// GetTemplate generates the requested template from the template map.
func GetTemplate(name string) *template.Template {
	return templatesF[name]()
}
