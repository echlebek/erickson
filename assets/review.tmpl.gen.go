// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package assets

var review = tmpl(asset.init(asset{Name: "review.tmpl", Content: "" +
	"<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\t{{ range $k, $v := .Stylesheets }}\n\t<link rel=\"stylesheet\" href=\"/assets/{{ $k }}\">\n\t{{ end }}\n\n\t{{ range $k, $v := .Scripts }}\n\t<script async=\"async\" src=\"/assets/{{ $k }}\"></script>\n\t{{ end }}\n  </head>\n  <body class=\"bg-base01\">\n\t<!-- the html in <aside> will be inserted into the diff table via jquery -->\n\t<aside class=\"hidden\">\n\t  <div class=\"bg-base2 base01 p2 border annotate-form\">\n\t    <pre class=\"line-text code-pre font-14\"></pre>\n\t    <br>\n\t    <br>\n\t\t<form method=\"post\" action=\"/reviews/{{.ID}}/annotations\">\n\t\t  {{ .CSRFField }}\n\t\t  <label for=\"comment\">Leave a comment</label>\n\t\t  <textarea name=\"comment\" id=\"comment\" class=\"monospace base01 bg-base3 block m1 pr3 field col-12\"></textarea>\n\t\t  <input id=\"file\" class=\"hidden\" name=\"file\"></input>\n\t\t  <input id=\"hunk\" class=\"hidden\" name=\"hunk\"></input>\n\t\t  <input id=\"line\" class=\"hidden\" name=\"line\"></input>\n\t\t  <button class=\"btn btn-outline\">Post Comment</button>\n\t\t  <div class=\"btn btn-outline\" onClick=\"toggleMenu();\">Cancel</div>\n\t    </form>\n\t  </div>\n\t</aside>\n\t<nav class=\"clearfix white bg-base02\">\n\t  <div class=\"sm-col\">\n\t\t<div class=\"menu-control\">\n\t\t\t<input class=\"hidden\" id=\"show-menu\" value=\"0\" name=\"\" type=\"checkbox\"></input>\n\t\t\t<label for=\"show-menu\" class=\"h1 ml2 aqua border-round\">&#9881;<span class=\"font-6\"> &#9660;</span></label>\n\t\t\t<a class=\"btn aqua mt1 mb1 ml1\" href=\"/\">Erickson Code Review</a>\n\t\t\t<a class=\"btn aqua mt1 mb1\" href=\"\">Review {{.ID}} ({{.Status}})</a>\n\t\t</div>\n\t  </div>\n\t  <div class=\"sm-col-right align-middle\">\n\t  {{ if .StatusOpen }}\n\t    <span>\n\t\t<form style=\"float: left;\" method=\"post\" action=\"/reviews/{{.ID}}/status\">\n\t\t  <button name=\"status\" value=\"Submitted\" class=\"btn btn-outline mt1 mb1 mr1 btn-submit\">Submit Review</button>\n\t\t  {{ .CSRFField }}\n\t\t</form>\n\t\t<form style=\"float: right;\" method=\"post\" action=\"/reviews/{{.ID}}/status\">\n\t\t  <button name=\"status\" value=\"Discarded\" class=\"btn btn-outline mt1 mb1 mr1 btn-discard\">Discard Review</button>\n\t\t  {{ .CSRFField }}\n\t\t</form>\n\t    </span>\n\t  {{ else }}\n\t\t<form method=\"post\" action=\"/reviews/{{.ID}}/status\">\n\t      <button name=\"status\" value=\"Open\" class=\"btn btn-outline mt1 mb1 mr1 btn-reopen\">Re-open Review</button>\n\t\t  {{ .CSRFField }}\n\t\t</form>\n\t  {{ end }}\n\t  </div>\n\t</nav>\n\t<div id=\"menu\" class=\"menu hidden clearfix bg-base02\">\n\t\t<div class=\"sm-col\">\n\t\t\t<p class=\"aqua mt2 ml2\">{{index .Session.Values \"username\" }}</p>\n\t\t\t<form action=\"/logout\" method=\"post\">\n\t\t\t\t{{ .CSRFField }}\n\t\t\t\t<button class=\"aqua ml2 mb2 btn btn-outline\">Log Out</button>\n\t\t\t</form>\n\t\t</div>\n\t</div>\n\t<div class=\"code-tables\">\n      {{ range $r, $revision := .R.Revisions }}\n\t  <a href=\"#/rev/{{ $r }}\" class=\"btn btn-outline base3\">Revision {{ $r }}</a>\n      {{ range $i, $f := $revision.Files }}\n\t  <table class=\"table table-condensed rounded-top mb1 mt1\" id=\"diff\">\n\t\t<tr class=\"rounded-top\">\n\t\t\t<th colspan=\"2\" class=\"rounded-top h5 bg-red\"><pre class=\"bold inline\"><code>--- {{ $f.OldName }}</code></pre></td>\n\t\t\t<th colspan=\"2\" class=\"rounded-top h5 bg-green\"><pre class=\"bold inline\"><code>+++ {{ $f.NewName }}</code></pre></td>\n\t\t</th>\n        {{ range $j, $hunk := $f.Hunks }}\n        {{ range $k, $lhs := $hunk.LHS }}\n\t\t<tr class=\"diffline\" id=\"diff-{{$i}}-{{$j}}-{{$k}}\" >\n\t\t  <td class=\"lineno lineno-lhs {{ (index $hunk.LHS $k).Type }}\">\n\t\t\t<pre class=\"code-pre h6\"><code> {{ (index $hunk.LHS $k).Line }}</code></pre>\n          </td>\n          <td class=\"line {{ (index $hunk.LHS $k).Type }}\">\n            <pre class=\"code-pre h6\"><code> {{ (index $hunk.LHS $k).Text }}</code></pre>\n\t\t  </td>\n\t\t  <td class=\"lineno lineno-rhs {{ (index $hunk.RHS $k).Type }}\">\n            <pre class=\"code-pre h6\"><code> {{ (index $hunk.RHS $k).Line }}</code></pre>\n          </td>\n          <td class=\"line {{ (index $hunk.RHS $k).Type }}\">\n            <pre class=\"code-pre h6\"><code> {{(index $hunk.RHS $k).Text }}</code></pre>\n          </td>\n\t  \t</tr>\n\t\t{{ if ($revision.GetAnnotations $i $j $k) }}\n\t\t<tr>\n\t\t  <td colspan=\"4\">\n\t\t\t<div class=\"bg-base2 px1 pt1 border\">\n\t\t\t  <p class=\"bold h4 regular mt1\">{{ $f.NewName }}:{{ (index $hunk.RHS $k).Line }}</p>\n\t\t\t  {{ range $a, $annotation := ($revision.GetAnnotations $i $j $k) }}\n\t\t\t  <p class=\"m1\">{{ $annotation.User }} says:</p>\n\t\t      <div class=\"bg-base3 p1 mb1 border rounded\">\n\t\t\t    <pre class=\"code-pre h6\"><code> {{ $annotation.Comment }}</code></pre>\n\t\t\t    <br>\n\t\t\t  </div>\n\t\t\t  {{ end }}\n\t\t\t</div>\n\t\t  </td>\n\t    </tr>\n\t\t{{ end }}\n        {{ end }}\n        <tr>\n          <td class=\"lineno\"></td>\n          <td class=\"line unchanged\"></td>\n          <td class=\"lineno\"></td>\n          <td class=\"line unchanged\"></td>\n        </tr>\n        {{ end }}\n      </table>\n      {{ end }}\n      {{ end }}\n    </div>\n  </div>\n</body>\n</html>\n" +
	""}))
