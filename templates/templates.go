// +build !dev

package templates

import (
	"html/template"
	"net/http"
)

//go:generate -command asset go run asset.go
//go:generate asset review.tmpl
//go:generate asset reviews.tmpl
//go:generate asset header.html

func tmpl(a asset) *template.Template {
	return template.Must(template.New(a.Name).Parse(a.Content))
}

type HTML template.HTML

func (h HTML) HTML() template.HTML {
	return template.HTML(h)
}

func html(a asset) HTML {
	return HTML(a.Content)
}

func css(a asset) http.Handler {
	return a
}

var (
	ReviewTmpl  = review
	ReviewsTmpl = reviews
	HeaderHTML  = header
)
