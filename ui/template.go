package ui

import (
	"html/template"
	"strings"
)

var (
	templatesF = map[string]func() *template.Template{
		"dashboard": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("single/dashboard.gohtml")))
			return b
		},
		"tasksAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("tasks/all.gohtml")))
			return b
		},
		"tasksSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("tasks/single.gohtml")))
			return b
		},
		"tasksNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("tasks/new.gohtml")))
			return b
		},
		"login": func() *template.Template {
			l := template.Must(template.New("login").Funcs(funcMap).Parse(fsMust("single/login.gohtml")))
			return l
		},
		"jobsAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("jobs/all.gohtml")))
			return b
		},
		"jobsSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("jobs/single.gohtml")))
			return b
		},
		"jobsNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("jobs/new.gohtml")))
			return b
		},
		"nodesAll": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("nodes/all.gohtml")))
			return b
		},
		"nodesSingle": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("nodes/single.gohtml")))
			return b
		},
		"nodesNew": func() *template.Template {
			b := baseTemplate()
			b = template.Must(b.Parse(fsMust("nodes/new.gohtml")))
			return b
		},
		"ec": func() *template.Template {
			t := template.Must(template.New("ec").Parse(fsMust("single/errorCode.gohtml")))
			return t
		},
		"form": func() *template.Template {
			t := template.Must(template.New("form").Funcs(funcMap).Parse(fsMust("single/formGen.gohtml")))
			return t
		},
	}
)

func baseTemplate() *template.Template {
	b := template.Must(template.New("base").Funcs(funcMap).Parse(fsMust("base/base.gohtml")))
	b = template.Must(b.Parse(fsMust("base/sidebar.gohtml")))
	b = template.Must(b.Parse(fsMust("base/header.gohtml")))
	b = template.Must(b.Parse(fsMust("base/footer.gohtml")))
	b = template.Must(b.Parse(fsMust("base/functions.gohtml")))
	return b
}

// Page describes the information that will be sent to the template.
type Page struct {
	Title   string
	URL     []string
	Details interface{}
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

var funcMap = template.FuncMap{
	"title":    strings.Title,
	"noescape": noescape,
}
