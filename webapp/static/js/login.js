$(document).ready(function(){
    $("button").click(function(){
        var useremail = $('#floatingInput').val()
        var password = $('#floatingPassword').val()

        if (!useremail) {

            return
        }else if (!password) {

            return
        }

        var data = {
            "useremail": useremail,
            "password": password
        }

        $.ajax({
            url: "/api/authorize",
            type: "POST",
            data: JSON.stringify(data),
            contentType: "application/json; charset=utf-8",
            success: function(data, textStatus, jqXHR){
                if (textStatus == "success"){
                    window.location.href = "/home";
                }
            },
            error: function(jqXHR, textStatus, errorThrown){
                if (textStatus == "success"){
                    window.location.href = "/home";
                }
            },
        });
    });
}); 