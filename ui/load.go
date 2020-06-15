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
		"base": func() *template.Template {
			b := template.Must(template.New("base").Parse(fsMust("/base.html")))
			return template.Must(b.Parse(fsMust("/sidebar.html")))
		}(),
	}
}
