$("#uploadForm").on("submit", upload);

function upload(e) {
  e.preventDefault();
  let image = $("#image").prop("files")[0];
  let shape = $("#shape").val();
  let numShapes = $("#numShapes").val();
  let formData = new FormData();
  formData.append("image", image);
  formData.append("shape", shape);
  formData.append("numShapes", numShapes);
  $.ajax({
    type: "POST",
    url: "http://localhost:8080/upload",
    data: formData,
    contentType: false,
    processData: false,
    beforeSend: function() {
      $("#result").removeAttr("hidden");
    },
    success: function(response) {
      console.log(response);
      $("#downloadResult").attr(
        "href",
        "https://res.cloudinary.com/robihid/image/upload/" + response
      );
      $("#downloadResult").removeAttr("hidden");
      $("#imgResult").attr(
        "src",
        "https://res.cloudinary.com/robihid/image/upload/" + response
      );
    }
  });
}
