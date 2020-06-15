// +build live

package ui

import "html/template"

func getTemplate(name string) *template.Template {
	return templatesF[name]()
}
