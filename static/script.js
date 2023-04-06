// JavaScript code
$(document).ready(function () {
    $('form').on('submit', function (event) {
        event.preventDefault();
        var formData = new FormData(this);
        var content = $('#text').val(); // get the value of the "content" input
        if (!content) {
            $('#desc').html('');
            $('#alert').html('Error: content cannot be empty.'); // display an error message if the "content" input is empty
            return;
        }
        $('#submit-btn').prop('disabled', true); // disable the button
        $('#text').val('');
        $('#alert').html('');
        $('#desc').html('');
        $('#spinner').show(); // show the spinner
        $.ajax({
            url: '/gen',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function (data) {
                dataObj = JSON.parse(data)
                if (dataObj.StatusCode === 0) {
                    $('#desc').html(content);
                    $('#image-block').attr("src", dataObj.Data)
                } else {
                    $('#alert').html(dataObj.Data);

                }
                $('#submit-btn').prop('disabled', false); // enable the button
                $('#spinner').hide(); // show the spinner

            }
        });
    });
});