function openNav() {
    document.getElementById("navbar").style.width = "248px";
    document.getElementById("content").style.marginLeft = "248px";
}

function closeNav() {
    document.getElementById("navbar").style.width = "0px";
    document.getElementById("content").style.marginLeft = "0px";
}

$("body").on("click","#nav-toggle", function(){
    var toggle = document.getElementById('nav-toggle');
    if (toggle.classList.contains("active")) {
        toggle.classList.remove("active");
        closeNav();
    } else {
        toggle.classList.add("active");
        openNav();
    }
});

function bringSettings() {
    var main = document.getElementById('navbar-contents-main');
    var settings = document.getElementById('navbar-contents-settings');
    main.hidden = true;
    settings.hidden = false;
}

function bringMain() {
    var main = document.getElementById('navbar-contents-main');
    var settings = document.getElementById('navbar-contents-settings');
    main.hidden = false;
    settings.hidden = true; 
}

$("body").on("click","#navbar-container-settings", function(){
    bringSettings();
});

$("body").on("click","#navbar-back-settings", function(){
    bringMain();
});

