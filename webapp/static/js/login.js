$(document).ready(function(){
    $("button").click(function(){
        $.post("/login",
        {
            email: "Donald Duck",
            password: "Duckburg"
        },
        function(data, status){
            console.log("Data: " + data + "\nStatus: " + status);
            window.location.href = "/home";
        });
    });
}); 