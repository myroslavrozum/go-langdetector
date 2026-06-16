$(document).ready(function() {
    $('#form').on('submit', function(event) {
        event.preventDefault(); 

        $.ajax({
            url: 'http://localhost:8080/detect',
            type: 'POST',
            data: $(this).serialize(),
            dataType: 'json',
            success: function(response) {
                alert('Form submitted successfully!');
                console.log(response);
            },
            error: function(xhr, status, error) {
                alert('An error occurred during submission.');
            }
        });
    });
});