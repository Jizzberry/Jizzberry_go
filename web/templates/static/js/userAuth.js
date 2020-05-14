var poll_usr;
var create_user = function(csrf_token, username, password, admin) {
    var txtUsername = username;
    var txtPassword = password;
    console.log(admin)
    poll_xhr = $.ajax({
        type: 'POST',
        url:   '/Jizzberry/settings/new_user/',
        data: {csrfmiddlewaretoken: csrf_token, username: txtUsername, password: txtPassword, manageServer: admin},
        success: function(result) {
        }
    });
}