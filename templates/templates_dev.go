// +build dev

package templates

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

//go:generate -command asset go run asset.go
//go:generate asset review.tmpl
//go:generate asset reviews.tmpl
//go:generate asset header.html

type Template struct {
	asset asset
}

// Wrap template.Template so we can conveniently reload the template
// from disk while developing.
func (t Template) Execute(w io.Writer, data interface{}) error {
	f, err := t.asset.Open()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New(t.asset.Name).Parse(string(content)))
	return tmpl.Execute(w, data)
}

func tmpl(a asset) Template {
	return Template{asset: a}
}

type HTML struct {
	asset asset
}

func (h HTML) HTML() template.HTML {
	f, err := h.asset.Open()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return template.HTML(content)
}

func html(a asset) HTML {
	return HTML{asset: a}
}

func css(a asset) http.Handler {
	return a
}

var (
	ReviewTmpl  = review
	ReviewsTmpl = reviews
	HeaderHTML  = header
)
