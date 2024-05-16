package templates

import "html/template"

type Templates struct {
	Templates *template.Template
}

func New() *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob("../../web/template/*.html")),
	}
}
