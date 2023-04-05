// JavaScript code
$(document).ready(function () {
    $('form').on('submit', function (event) {
        event.preventDefault();
        var formData = new FormData(this);
        var content = $('#text').val(); // get the value of the "content" input
        if (!content) {
            $('#alert').attr("class", "alert alert-danger")
            $('#alert').show();
            $('#content-label').html('Error: content cannot be empty.'); // display an error message if the "content" input is empty
            return;
        }
        $('#submit-btn').prop('disabled', true); // disable the button
        $('#text').val('');
        $('#alert').attr("class", "alert alert-primary")
        $('#alert').show();
        $('#content-label').html(content);
        $('#spinner').show(); // show the spinner
        $.ajax({
            url: '/gen',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function (data) {
                
                $('#image-block').html('<img class="w-100" src="' + data + '">');
                $('#spinner').hide(); // hide the spinner after receiving the response
                $('#submit-btn').prop('disabled', false); // enable the button
            }
        });
    });
});