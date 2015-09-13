package assets

import "net/http"

//go:generate -command asset go run asset.go
//go:generate asset review.tmpl
//go:generate asset reviews.tmpl
//go:generate asset styles.css
//go:generate asset basscss.css
//go:generate asset app.js
//go:generate asset jquery.js

var (
	ScriptHandlers = map[string]http.Handler{
		"jquery.js": jquery,
		"app.js":    app,
	}
	StylesheetHandlers = map[string]http.Handler{
		"styles.css":  styles,
		"basscss.css": basscss,
	}
	Templates = map[string]Template{
		"review.html":  review,
		"reviews.html": reviews,
	}
)

func css(a asset) asset {
	return a
}

func js(a asset) asset {
	return a
}
