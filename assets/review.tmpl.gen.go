// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package assets

var review = tmpl(asset.init(asset{Name: "review.tmpl", Content: "" +
	"<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\t{{ range $k, $v := .Stylesheets }}\n\t<link rel=\"stylesheet\" href=\"/assets/{{ $k }}\">\n\t{{ end }}\n\n\t{{ range $k, $v := .Scripts }}\n\t<script src=\"/assets/{{ $k }}\"></script>\n\t{{ end }}\n  </head>\n  <body>\n  \t<nav class=\"clearfix white bg-black\">\n\t  <div class=\"sm-col\">\n\t\t<a class=\"btn aqua mt1 mb1 ml1\" href=\"/\">Erickson Code Review</a>\n\t\t<a class=\"btn aqua mt1 mb1\" href=\"\">Review {{.ID}} ({{.Status}})</a>\n\t  </div>\n\t  <div class=\"sm-col-right align-middle\">\n\t  {{ if .StatusOpen }}\n\t    <button onClick=\"submitReview()\" class=\"btn btn-outline mt1 mb1 mr1 btn-submit\">Submit Review</button>\n\t    <button onClick=\"discardReview()\" class=\"btn btn-outline mt1 mb1 mr1 btn-discard\">Discard Review</button>\n\t  {{ else }}\n\t    <button onClick=\"reopenReview()\" class=\"btn btn-outline mt1 mb1 mr1 btn-reopen\">Re-open Review</button>\n\t  {{ end }}\n\t  </div>\n\t</nav>\n\t<div class=\"code-tables\">\n      {{ range $r, $revision := .R.Revisions }}\n      <h3>Revision {{ $r }}</h3>\n      {{ range $i, $f := $revision.Files }}\n      <table style=\"width: 100%; display: table; border-width: thin; border-style: solid;\" class=\"table table-condensed mb1\" id=\"diff\">\n        <td class=\"lineno\"></td>\n        <td class=\"fileheader\"><pre class=\"blue code-pre\"><code>--- {{ $f.OldName }}</code></pre></td>\n        <td class=\"fileheader\"><pre class=\"blue code-pre\"><code>+++ {{ $f.NewName }}</code></pre></td>\n        <td class=\"lineno\"></td>\n        {{ range $j, $hunk := $f.Hunks }}\n        {{ range $k, $lhs := $hunk.LHS }}\n        <tr style=\"padding: 0px; border: none\">\n          <td class=\"lineno {{ (index $hunk.LHS $k).Type }}\">\n            <pre class=\"code-pre\"><code> {{ (index $hunk.LHS $k).Line }}</code></pre>\n          </td>\n          <td class=\"line {{ (index $hunk.LHS $k).Type }}\">\n            <pre class=\"code-pre\"><code> {{ (index $hunk.LHS $k).Text }}</code></pre>\n          </td>\n          <td class=\"line {{ (index $hunk.RHS $k).Type }}\">\n            <pre class=\"code-pre\"><code> {{(index $hunk.RHS $k).Text }}</code></pre>\n          </td>\n          <td class=\"lineno {{ (index $hunk.RHS $k).Type }}\">\n            <pre class=\"code-pre\"><code> {{ (index $hunk.RHS $k).Line }}</code></pre>\n          </td>\n        </tr>\n        {{ end }}\n        <tr>\n          <td class=\"lineno\"><pre class=\"code-pre\"><code>  &#10714;</code></pre></td>\n          <td class=\"line unchanged\"></td>\n          <td class=\"line unchanged\"></td>\n          <td class=\"lineno\"><pre class=\"code-pre\"><code>  &#10715;</code></pre></td>\n        </tr>\n        {{ end }}\n      </table>\n      {{ end }}\n      {{ end }}\n    </div>\n  </div>\n</body>\n</html>\n" +
	""}))
