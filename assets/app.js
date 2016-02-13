function createReview() {
  $("#create-review").removeClass("hidden");
  $("#cancel-button").removeClass("hidden");
}

function cancelReview() {
  $("#create-review").addClass("hidden");
  $("#cancel-button").addClass("hidden");
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

function patchRevision(annotation) {
  var data = JSON.stringify(annotation);
  console.log(data);
  $.ajax({
    headers: {
      "Content-Type": "application/json"
    },
    url: window.location,
    type: "PATCH",
    data: data,
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
  $("#diffcheckmark").show();
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

  $("#show-menu").change(function () {
    if (this.checked) {
      $("#menu").show();
    } else {
      $("#menu").hide(); 
    }
  });
}

function showAnnotate(span) {
  cancelAnnotate(); // If another annotate dialog is visible, remove it

  // Add the annotate form after the selected row
  var form = $(".annotate-form");
  form.show();
  var div = $(span).parent();
  div.after(form);
  form.wrap('<div></div>');

  // Put some descriptive information into the annotate form
  var spans = $(span).parent().children("span");
  var lineTd = $(spans[2]);
  var diffTd = $(spans[3]);
  var diffText = diffTd.find("code").text();
  var lineText = lineTd.find("code").text();
  form.find(".line-text").text(lineText + ": " + diffText);

  // Stash some data in the comment element for submitting to the server
  var idParts = div.attr("id").split("-");
  $("#file").val(Number(idParts[1]));
  $("#hunk").val(Number(idParts[2]));
  $("#line").val(Number(idParts[3]));
}

function cancelAnnotate() {
  $(".annotate-form").hide();
}
