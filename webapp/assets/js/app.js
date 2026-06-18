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

    //let socket = new WebSocket("wss://javascript.info/article/websocket/demo/hello");
    let socket = new WebSocket("ws://localhost:8080/echo");
    socket.onopen = function(e) {
        console.log("[open] Connection established");
        console.log("Sending to server");
        socket.send("My name is John");
    };

    socket.onmessage = function(event) {
        console.log(`[message] Data received from server: ${event.data}`);
        $('#logger').append(`${event.data}<br>`);
    };

    socket.onclose = function(event) {
        if (event.wasClean) {
            console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
        // e.g. server process killed or network down
        // event.code is usually 1006 in this case
            console.log('[close] Connection died');
        }
    };

    socket.onerror = function(error) {
        console.log(`[error]`);
    };
});