// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package assets

var app = js(asset.init(asset{Name: "app.js", Content: "" +
	"function createReview() {\n  $(\"#create-review\").removeClass(\"hidden\");\n  $(\"#cancel-button\").removeClass(\"hidden\");\n}\n\nfunction cancelReview() {\n  $(\"#create-review\").addClass(\"hidden\");\n  $(\"#cancel-button\").addClass(\"hidden\");\n}\n\nfunction submitReview() {\n  patchReview(\"Submitted\");\n}\n\nfunction discardReview() {\n  patchReview(\"Discarded\");\n}\n\nfunction reopenReview() {\n  patchReview(\"Open\");\n}\n\nfunction toggleShowAll() {\n  if ($(\"#show-all\").is(\":checked\")) {\n    $(\"tr\").show();\n  } else {\n    $(\"td > span\").filter(function () {\n      return $(this).text() !== \"Open\";\n    }).parent().parent().hide();\n  }\n}\n\nfunction patchReview(status) {\n  $.ajax({\n    headers: {\n      \"Content-Type\": \"application/json\"\n    },\n    url: window.location,\n    type: \"PATCH\",\n    data: JSON.stringify({status: status}),\n    complete: function() {\n      window.location.reload();\n    }\n  });\n}\n\nfunction patchRevision(annotation) {\n  $.ajax({\n    headers: {\n      \"Content-Type\": \"application/json\"\n    },\n    url: window.location,\n    type: \"PATCH\",\n    data: JSON.stringify(annotation),\n    complete: function() {\n      window.location.reload();\n    }\n  });\n}\n\nfunction pasteFile(file) {\n  var reader = new FileReader();\n  reader.onload = function(e) {\n    $(\"#diff\").text(e.target.result);\n  }\n  reader.readAsText(file);\n}\n\nfunction annotate(fileName, revision, hunk, line) {\n  var req = {\n    fileName: fileName,\n    hunk: hunk,\n    line: line,\n    message: message\n  };\n  console.log(req);\n  patchRevision(req);\n}\n\nwindow.onload = function() {\n  // Show only the selected reviews\n  toggleShowAll();\n  $(document).on(\"click\", \".lineno-lhs\", function () {\n    showAnnotate(this);\n  });\n\n  $(document).on(\"click\", \".lineno-rhs\", function () {\n    showAnnotate(this);\n  });\n}\n\nfunction showAnnotate(td) {\n  cancelAnnotate(); // If another annotate dialog is visible, remove it\n\n  // Add the annotate form after the selected row\n  var form = $(\".annotate-form\");\n  form.show();\n  var tr = $(td).parent();\n  tr.after(form);\n  form.wrap('<tr><td colspan=\"4\"></td></tr>');\n\n  // Put some descriptive information into the annotate form\n  var tds = $(td).parent().children(\"td\");\n  var lineTd = $(tds[2]);\n  var diffTd = $(tds[3]);\n  var diffText = diffTd.find(\"code\").text();\n  var lineText = lineTd.find(\"code\").text();\n  form.find(\".line-text\").text(lineText + \": \" + diffText);\n\n  // Stash some data in the comment element for submitting to the server\n  var idParts = tr.attr(\"id\").split(\"-\");\n  $(\"#comment\").data(\"fileNumber\", int(idParts[1]));\n  $(\"#comment\").data(\"hunkNumber\", int(idParts[2]));\n  $(\"#comment\").data(\"lineNumber\", int(idParts[3]));\n}\n\nfunction cancelAnnotate() {\n  $(\".annotate-form\").hide();\n}\n\nfunction postComment() {\n  var comment = $(\"#comment\");\n  var annotation = {\n    fileNumber: comment.data(\"fileNumber\"),\n    hunkNumber: comment.data(\"hunkNumber\"),\n    lineNumber: comment.data(\"lineNumber\"),\n    comment: comment.text()\n  };\n  patchRevision(annotation);\n}\n" +
	""}))