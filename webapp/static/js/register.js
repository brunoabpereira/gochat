$(document).ready(function(){
    var successHTML = `
        <h1 class="display-5 fw-bold bs-primary-text-emphasis" style="text-align: center">
            <a href="/" class="text-decoration-none" >gochat</a>  
        </h1>
        <h3 class="bs-primary-text-emphasis" style="text-align: center">
            Success!
        </h3>
        <h5 class="bs-primary-text-success">
            Your account has been created.
        </h5>
        <a class="btn btn-success w-100 py-2" href="/home">Continue</a>
    `
    $("button").click(function(){
        $.post("/login",
        {
            email: "Donald Duck",
            password: "Duckburg"
        },
        function(data, status){
            console.log("Data: " + data + "\nStatus: " + status);
            $('#regForm').fadeOut(100, function() {
                $(this).html(successHTML).fadeIn(100);
            });
        });
    });
}); 