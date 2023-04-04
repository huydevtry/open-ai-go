
$(document).ready(function() {
    $('form').on('submit', function(event) {
      event.preventDefault();
      var formData = new FormData(this);
      $.ajax({
        url: '/upload',
        type: 'POST',
        data: formData,
        processData: false,
        contentType: false,
        success: function(data) {
          console.log(data)
          const jdata = JSON.parse(data)
          $('#image-block').html('<p>' + jdata.data[0].b64_json + '</p>');
        }
      });
    });
  });
  