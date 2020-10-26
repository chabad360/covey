package ui

import (
	"html/template"
	"strings"
)

var (
	templatesF = map[string]func() *template.Template{
		"dashboard": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/dashboard.html")))
			return b
		},
		"tasksAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/tasks/all.html")))
			return b
		},
		"tasksSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/tasks/single.html")))
			return b
		},
		"tasksNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/tasks/new.html")))
			return b
		},
		"login": func() *template.Template {
			l := template.Must(template.New("login").Funcs(funcMap).Parse(fsMust("/single/login.html")))
			return l
		},
		"jobsAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/jobs/all.html")))
			return b
		},
		"jobsSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/jobs/single.html")))
			return b
		},
		"jobsNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/jobs/new.html")))
			return b
		},
		"nodesAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/nodes/all.html")))
			return b
		},
		"nodesSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/nodes/single.html")))
			return b
		},
		"nodesNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("/nodes/new.html")))
			return b
		},
		"ec": func() *template.Template {
			t := template.Must(template.New("ec").Parse(fsMust("/single/errorCode.html")))
			return t
		},
		"pluginForm": func() *template.Template {
			t := template.Must(template.New("form").Parse(fsMust("/single/formGen.html")))
			return t
		},
	}
)

func baseTemplate() *template.Template {
	b := template.Must(template.New("base").Funcs(funcMap).Parse(fsMust("/base/base.html")))
	b = template.Must(b.Parse(fsMust("/base/sidebar.html")))
	b = template.Must(b.Parse(fsMust("/base/header.html")))
	b = template.Must(b.Parse(fsMust("/base/footer.html")))
	b = template.Must(b.Parse(fsMust("/base/functions.html")))
	return b
}

// Page describes the information that will be sent to the template.
type Page struct {
	Title   string
	URL     []string
	Details interface{}
}

var funcMap = template.FuncMap{
	"title": strings.Title,
}
