function openNav() {
    document.getElementById("sidebar").style.width = "13rem";
    document.getElementById("home-nav").style.width = "13rem";
    document.getElementById("settings-nav").style.width = "13rem";
    document.getElementById("content").style.paddingLeft = "13rem";
}

function closeNav() {
    document.getElementById("sidebar").style.width = "0rem";
    document.getElementById("home-nav").style.width = "0rem";
    document.getElementById("settings-nav").style.width = "0rem";
    document.getElementById("content").style.paddingLeft = "0rem";
}

window.onpaint = navbarSelector();

function navbarSelector(){
    var currentUrl = window.location.pathname;
    if (currentUrl.toLowerCase() == "/jizzberry/settings") {
        handleSettingsPress()
    }
}

$("body").on("click","#nav-toggle", function(){
    var toggle = document.getElementById('nav-toggle');
    if (toggle.classList.contains("active")) {
        console.log("closing")
        toggle.classList.remove("active");
        closeNav();
    } else {
        toggle.classList.add("active");
        console.log("opening")
        openNav();
    }
});

function handleTabChange(name) {
	var nodes = document.getElementById("content").childNodes;
	
	for (var i=0; i<nodes.length; i++) {
		if (nodes[i].nodeName.toLowerCase() == 'div') {
			var child = nodes[i];
			child.hidden = true;
		}
	}
	
	document.getElementById(name).hidden = false;
}

function handleBackPress() {
	document.getElementById("settings-nav").hidden = true;
	document.getElementById("home-nav").hidden = false;
}

function handleSettingsPress() {
    var currentUrl = window.location.pathname;

    if (currentUrl.toLowerCase() == "/jizzberry/settings") {
        document.getElementById("settings-nav").hidden = false;
        document.getElementById("home-nav").hidden = true;
    } else {
        window.location.href = '/Jizzberry/settings';
    }
}