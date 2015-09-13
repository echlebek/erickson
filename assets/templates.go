// +build !dev

package assets

import (
	"html/template"
	"io"
)

type Template struct {
	asset asset
}

// Wrap template.Template so we can conveniently reload the template
// from disk while developing. Embeds any scripts or stylesheets.
func (t Template) Execute(w io.Writer, data interface{}) error {
	tmpl := template.Must(template.New(t.asset.Name).Parse(t.asset.Content))
	return tmpl.Execute(w, data)
}

func tmpl(a asset) Template {
	return Template{asset: a}
}
