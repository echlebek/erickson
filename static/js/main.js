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

function pasteFile(file) {
  var reader = new FileReader();
  reader.onload = function(e) {
    $("#diff").text(e.target.result);
  }
  reader.readAsText(file);
}

window.onload = toggleShowAll;
