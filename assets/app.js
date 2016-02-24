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

function publishAnnotations() {
  $.ajax({
    headers: {
      "Content-Type": "application/json"
    },
    url: window.location + "/annotations",
    type: "PATCH",
    data: "",
    complete: function () {
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

function keyEvent(event) {
  // 44.4444 is roughly the amount chrome scrolls with one down-arrow key
  // press on my macbook pro.
  var chr = event.char;
  if (!chr) { // web browsers...
    chr = event.keyCode == 106 ? "j" : (event.keyCode == 107 ? "k" : null);
  }
  if (chr == "j") {
    window.scrollTo(0, window.scrollY + 44.4444);
  } else if (chr == "k") {
    window.scrollTo(0, window.scrollY - 44.4444);
  }
}
