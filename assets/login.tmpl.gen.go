// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package assets

var login = tmpl(asset.init(asset{Name: "login.tmpl", Content: "" +
	"<!DOCTYPE html>\n<html lang=\"en\">\n\t<head>\n\t\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\t\t{{ range $k, $v := .Stylesheets }}\n\t\t<link rel=\"stylesheet\" href=\"/assets/{{ $k }}\">\n\t\t{{ end }}\n\n\t\t{{ range $k, $v := .Scripts }}\n\t\t<script async=\"async\" src=\"/assets/{{ $k }}\"></script>\n\t\t{{ end }}\n\t</head>\n\t<body class=\"bg-base01\">\n\t\t<nav class=\"clearfix white bg-base02\">\n\t\t<div class=\"sm-col\">\n\t\t\t<a class=\"btn aqua mt1 mb1 ml1\" href=\"/\">Erickson Code Review</a>\n\t\t\t<a class=\"btn aqua mt1 mb1\" href=\"\">Login</a>\n\t\t</div>\n\t\t</nav>\n\t\t<div class=\"container base02 bg-base3 border rounded mt4 p2\">\n\t\t\t<h2 class=\"h2\">Log In</h2>\n\t\t\t<form method=\"post\" action=\"/login\">\n\t\t\t\t<div class=\"p2\">\n\t\t\t\t\t{{ .CSRFField }}\n\t\t\t\t\t<input class=\"m1 block col-12 field\" type=\"email\" name=\"username\" placeholder=\"Username\">\n\t\t\t\t\t<input class=\"m1 block col-12 field\" type=\"password\" name=\"password\" placeholder=\"Password\">\n\t\t\t\t\t<button type=\"submit\" class=\"m1 btn btn-outline\">Submit</button> \n\t\t\t\t</div>\n\t\t\t</form>\n\t\t</div>\n\t</body>\n</html>\n" +
	""}))
