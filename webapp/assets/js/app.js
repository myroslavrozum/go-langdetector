$(document).ready(function() {
    $('#form').on('submit', function(event) {
        event.preventDefault(); 

        $.ajax({
            url: 'http://localhost:8080/detect',
            type: 'POST',
            data: $(this).serialize(),
            dataType: 'json',
            success: function(response) {
                console.log('Form submitted successfully, detected language is: '+ response.minLangFull);
                $("*").removeClass("active")
                $('#langmarker_' + response.minLang).addClass('active');
            },
            error: function(xhr, status, error) {
                alert('An error occurred during submission.');
            }
        });
    });
});