package templates

import "html/template"

//go:generate -command asset go run asset.go
//go:generate asset review.tmpl
//go:generate asset reviews.tmpl
//go:generate asset header.html

func tmpl(a asset) *template.Template {
	return template.Must(template.New(a.Name).Parse(a.Content))
}

func html(a asset) template.HTML {
	return template.HTML(a.Content)
}

var (
	ReviewTmpl  = review
	ReviewsTmpl = reviews
	HeaderHTML  = header
)
