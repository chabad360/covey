package ui

import (
	"html/template"
	"strings"
)

var (
	templatesF = map[string]func() *template.Template{
		"base": func() *template.Template {
			b := template.Must(template.New("base").Funcs(funcMap).Parse(fsMust("/base.html")))
			b = template.Must(b.Parse(fsMust("/sidebar.html")))
			b = template.Must(b.Parse(fsMust("/header.html")))
			b = template.Must(b.Parse(fsMust("/footer.html")))
			return b
		},
		"login": func() *template.Template {
			l := template.Must(template.New("login").Funcs(funcMap).Parse(fsMust("/login.html")))
			return l
		},
	}
)

type page struct {
	Title string
	URL   []string
}

var funcMap = template.FuncMap{
	"title": strings.Title,
}
