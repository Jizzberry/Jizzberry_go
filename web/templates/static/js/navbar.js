function toggleNav() {
    document.getElementById("sidebar").classList.toggle("toggled")
    document.getElementById("home-nav").classList.toggle("toggled")
    if (document.getElementById("settings-nav") != null) {
        document.getElementById("settings-nav").classList.toggle("toggled")
    }
    document.getElementById("content").classList.toggle("content-toggled")
}

window.onpaint = getPage();

function handleTabChange(name, element) {
    var nodes = document.getElementById("content").childNodes;

    for (var i = 0; i < nodes.length; i++) {
        if (nodes[i].nodeName.toLowerCase() === 'div') {
            var child = nodes[i];
            child.hidden = true;
        }
    }

    document.getElementById(name).hidden = false;

    setActiveSettingsButton(element)
}

function handleBackPress() {
	document.getElementById("settings-nav").hidden = true;
	document.getElementById("home-nav").hidden = false;
}

function handleSettingsPress() {
    var currentUrl = window.location.pathname;

    if (currentUrl.toLowerCase() === "/jizzberry/settings") {
        document.getElementById("settings-nav").hidden = false;
        document.getElementById("home-nav").hidden = true;
    } else {
        window.location.href = '/Jizzberry/settings';
    }
}

function setActiveSettingsButton(element) {
    let container = document.getElementById("settings-nav");
    for (var i = 0; i < container.childElementCount; i++) {
        if (container.children[i].classList.contains("active")) {
            container.children[i].classList.remove("active")
        }
    }

    element.classList.add("active")
}

function getPage() {
    let url = window.location.pathname;
    let n = url.lastIndexOf('/');
    if (n >= 0) {
        url = url.substring(n + 1);
    }

    switch (url.toLowerCase()) {
        case "home":
            setActiveButton("home");
            break;
        case "actors":
            setActiveButton("actors");
            break;
        case "tags":
            setActiveButton("tags");
            break;
        case "settings":
            setActiveButton("settings");
            handleSettingsPress()
            break;
    }
}

function setActiveButton(id) {
    let container = document.getElementById("home-nav");
    for (var i = 0; i < container.childElementCount; i++) {
        if (container.children[i].classList.contains("active")) {
            container.children[i].classList.remove("active");
        }
    }
    document.getElementById(id).classList.add("active");
}