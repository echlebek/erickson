function createReview() {
  $("#create-review").removeClass("hidden");
  $("#cancel-button").removeClass("hidden");
}

function cancelReview() {
  $("#create-review").addClass("hidden");
  $("#cancel-button").addClass("hidden");
}

function submitReview() {
  patchReview("Submitted");
}

function discardReview() {
  patchReview("Discarded");
}

function reopenReview() {
  patchReview("Open");
}

function toggleShowAll() {
  if ($("#show-all").is(":checked")) {
    $("tr").show();
  } else {
    $("td > span").filter(function () {
      return $(this).text() !== "Open";
    }).parent().parent().hide();
  }
}

function patchReview(status) {
  $.ajax({
    headers: {
      "Content-Type": "application/json"
    },
    url: window.location,
    type: "PATCH",
    data: JSON.stringify({status: status}),
    complete: function() {
      window.location.reload();
    }
  });
}

function patchRevision(annotation) {
  $.ajax({
    headers: {
      "Content-Type": "application/json"
    },
    url: window.location,
    type: "PATCH",
    data: JSON.stringify(annotation),
    complete: function() {
      window.location.reload();
    }
  });
}

function pasteFile(file) {
  var reader = new FileReader();
  reader.onload = function(e) {
    $("#diff").text(e.target.result);
  }
  reader.readAsText(file);
}

function annotate(fileName, revision, hunk, line) {
  var req = {
    fileName: fileName,
    hunk: hunk,
    line: line,
    message: message
  };
  console.log(req);
  patchRevision(req);
}

window.onload = function() {
  // Show only the selected reviews
  toggleShowAll();
  $(document).on("click", ".lineno-lhs", function () {
    showAnnotate(this);
  });

  $(document).on("click", ".lineno-rhs", function () {
    showAnnotate(this);
  });
}

function showAnnotate(td) {
  cancelAnnotate(); // If another annotate dialog is visible, remove it

  // Add the annotate form after the selected row
  var form = $(".annotate-form");
  form.show();
  var tr = $(td).parent();
  tr.after(form);
  form.wrap('<tr><td colspan="4"></td></tr>');

  // Put some descriptive information into the annotate form
  var tds = $(td).parent().children("td");
  var lineTd = $(tds[2]);
  var diffTd = $(tds[3]);
  var diffText = diffTd.find("code").text();
  var lineText = lineTd.find("code").text();
  form.find(".line-text").text(lineText + ": " + diffText);

  // Stash some data in the comment element for submitting to the server
  var idParts = tr.attr("id").split("-");
  $("#comment").data("fileNumber", int(idParts[1]));
  $("#comment").data("hunkNumber", int(idParts[2]));
  $("#comment").data("lineNumber", int(idParts[3]));
}

function cancelAnnotate() {
  $(".annotate-form").hide();
}

function postComment() {
  var comment = $("#comment");
  var annotation = {
    fileNumber: comment.data("fileNumber"),
    hunkNumber: comment.data("hunkNumber"),
    lineNumber: comment.data("lineNumber"),
    comment: comment.text()
  };
  patchRevision(annotation);
}
