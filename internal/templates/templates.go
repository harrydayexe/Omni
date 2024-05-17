package templates

import "html/template"

type Templates struct {
	Templates *template.Template
}

func New(templateGlob string) *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob(templateGlob)),
	}
}
