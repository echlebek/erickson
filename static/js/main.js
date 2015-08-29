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

function patchReview(status) {
  $.ajax({
    headers: {
      "Content-Type": "application/json"
    },
    url: window.location,
    type: "PATCH",
    data: JSON.stringify({status: status})
  });
}

function pasteFile(file) {
  var reader = new FileReader();
  reader.onload = function(e) {
    $("#diff").text(e.target.result);
  }
  reader.readAsText(file);
}
