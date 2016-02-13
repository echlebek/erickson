// +build dev

package assets

import (
	"html/template"
	"io"
	"io/ioutil"
)

type Template struct {
	asset asset
}

// Wrap template.Template to embed scripts and stylesheets.
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
