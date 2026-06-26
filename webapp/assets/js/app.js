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
                $("*").removeClass("detected")
                $('#langmarker_' + response.minLang).addClass('detected');
            },
            error: function(xhr, status, error) {
                alert('An error occurred during submission.');
            }
        });
    });

    //let socket = new WebSocket("wss://javascript.info/article/websocket/demo/hello");
    let logs = new EventSource("/logStream");
   
    logs.onmessage = function(event) {
        console.log("[message] Data received from server: ${event.data}");
        $('#logger').append(`${event.data}<br>`);
    };

    logs.onerror = function(error) {
        console.error("[error]", error);
    };
});