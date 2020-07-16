$(function () {
	$("#usercreationform").submit(function (event) {
		event.preventDefault();
		$.ajax({
			type: "POST",
			url: "/auth/create/",
			data: $(this).serialize(),
            success: function () {
                window.location.reload();
            },
            error: function (msg) {
                alert(msg);
            },
        });
    });
});

const clickButton = document.querySelector(".clickBtn");
const closeButton = document.querySelector(".close");
const popup = document.querySelector(".bottom-up");

clickButton.addEventListener('click', (e) => {
    popup.classList.add("show");
});

closeButton.addEventListener('click', (e) => {
    popup.classList.remove("show");
});

function postConfig() {
    const inputFolder = document.getElementById("folderForm");
    const inputFile = document.getElementById("fileForm");

    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
		}
	};

	xhr.open("POST", "/api/config", true);
	xhr.send(
		JSON.stringify({
			file_rename_formatter: inputFile.value,
			folder_rename_formatter: inputFolder.value,
		})
	);
}

function removePath(path) {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;

	xhr.onreadystatechange = function () {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			window.location.reload();
		}
	};

	xhr.open("DELETE", "/api/setPath?path=" + path, true);
	xhr.send();
}

function startScan() {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;
	xhr.open("POST", "/api/startScanTask", true);
	xhr.send();
}

function stopTask(name) {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;
	xhr.open("POST", "/api/stopTask?uid=" + name, true);
	xhr.send();

	resetControls(name);
}

const sessionsSocket = new WebSocket(
	"ws://" + window.location.host + "/ws/session" + ""
);

sessionsSocket.onmessage = function (e) {
	let parsed = JSON.parse(e.data);
	if (parsed["type"] === "progress") {
        console.log(parsed)
        handleProgress(parsed);
    } else if (parsed["type"] === "status") {
		setUserStatus(parsed);
	}
};

sessionsSocket.onclose = function (e) {
	console.error(e);
};

function handleProgress(data) {
    let value = data["value"]
    let text = document.getElementById(value["uid"] + "-textprogress");
    let bar = document.getElementById(value["uid"] + "-progressbar");

    if (text == null || bar == null) {
        createTask(value["name"], value["uid"], value["progress"]);
    } else {
        text.textContent = value["progress"] + "%";
        text.hidden = false;
        bar.style.width = value["progress"] + "%";
    }
}

function handleDropdown(name) {
	let dropdown = document.getElementById(name + "-angle-down");
	let task = document.getElementById(name + "-body");

	if (dropdown.classList.contains("active")) {
		task.style.height = "0";
		dropdown.classList.remove("active");
	} else {
		task.style.height = "initial";
		dropdown.classList.add("active");
	}
}

function dropdownSetter(name) {
	let dropdown = document.getElementById(name + "-angle-down");

	if (dropdown != null) {
		dropdown.setAttribute("onclick", 'handleDropdown("' + name + '");');
	}
}

function createTask(name, uid, progress) {
    let container = document.getElementById("tasks-container");

    let outermostDiv = document.createElement("div");
    outermostDiv.className = "task-stats mb-4";

    let topTextHolder = document.createElement("div");
    topTextHolder.className = "stats-head d-flex";
    outermostDiv.appendChild(topTextHolder);

    let flexDiv = document.createElement("div");
	flexDiv.className = "d-flex ml-2";
	topTextHolder.appendChild(flexDiv);

    let emptySpan = document.createElement("span");

    let paragraphElement0 = document.createElement("p");
    paragraphElement0.style.fontWeight = "700";
    paragraphElement0.className = "text-warning mr-3 mb-0 pt-2";
    paragraphElement0.textContent = "30 mins";

    let paragraphElement1 = document.createElement("p");
    paragraphElement1.style.fontWeight = "700";
    paragraphElement1.className = "mr-3 mb-0 pt-2";
    paragraphElement1.id = uid + "-textprogress";
    paragraphElement1.textContent = progress + "%";

    let paragraphElement2 = document.createElement("p");
    paragraphElement2.style.fontStyle = "italic";
    paragraphElement2.className = "ml-3 mb-0 pt-2";
    paragraphElement2.textContent = name;

    flexDiv.appendChild(paragraphElement0);
    flexDiv.appendChild(paragraphElement1);
    flexDiv.appendChild(paragraphElement2);

    let buttonsDiv = document.createElement("div");
    buttonsDiv.className = "stats-btn-grp d-flex";
    topTextHolder.appendChild(buttonsDiv);

    let makeButton = function (name, type, callback) {
        let buttonElement = document.createElement("button");
        buttonElement.style.width = "3rem";
        buttonElement.className = "btn";
        buttonElement.id = name + "-" + type;
        buttonElement.onclick = callback

        let icon = document.createElement("i");
        icon.className = "fas fa-" + type;

        buttonElement.appendChild(icon);
        return buttonElement;
    };

    buttonsDiv.appendChild(makeButton(uid, "stop", function () {
        stopTask(uid)
    }));
    buttonsDiv.appendChild(makeButton(uid, "angle-down"));

    let taskBody = document.createElement("div");
    taskBody.className = "stats-body d-flex flex-column";
    taskBody.style.height = "0";
    taskBody.style.overflow = "hidden";
    taskBody.id = uid + "-body";
    outermostDiv.appendChild(taskBody);

    let divProgressContainer = document.createElement("div");
    divProgressContainer.className = "progress my-3";

    let progressbar = document.createElement("div");
    progressbar.className = "progress-bar";
    progressbar.role = "progressbar";
    progressbar.id = uid + "-progressbar";

    divProgressContainer.appendChild(progressbar);

    let divSubText = document.createElement("div");
    divSubText.className = "p-grp d-flex";
    divSubText.innerHTML =
        '<p class="mr-5"><strong class="mr-1">Started:</strong><span>69 min ago</span></p><p class="mr-5"><strong class="mr-1">Status:</strong><span>alive</span></p>';

    taskBody.appendChild(divProgressContainer);
    taskBody.appendChild(divSubText);

	container.appendChild(outermostDiv);
}

function reqListener() {
	const parser = new DOMParser();
	let htmlDoc = parser.parseFromString(this.responseText, "text/html");
	let files = htmlDoc.getElementsByTagName("a");
	for (let i = 0; i < files.length; i++) {
		console.log(files[i].href);
		createLogTab(files[i].textContent.split(".")[0].replace(" ", "-"));
		createLogContent(
			files[i].textContent.replace(" ", "-").split(".")[0],
			"/logs/" + files[i].textContent.replace(" ", "%20")
		);
	}
}

function createLogTab(name) {
	let container = document.getElementById("list-tab");
	let fileTab = document.createElement("a");
	fileTab.className = "list-group-item list-group-item-action";
	fileTab.href = "#" + name;
	fileTab.setAttribute("data-toggle", "list");
	fileTab.role = "tab";
	fileTab.textContent = name;

	container.appendChild(fileTab);
}

function createLogContent(name, url) {
	function logsContentHandler() {
		let container = document.getElementById("nav-tabContent");
		let logsContent = document.createElement("div");
		logsContent.className = "tab-pane fade p-2";
		logsContent.style.boxShadow = "0 0 10px 1px lightgray";
		logsContent.style.borderRadius = "7px";
		logsContent.style.height = "30rem";
		logsContent.style.overflow = "scroll";
		logsContent.id = name;
		logsContent.role = "tabpanel";

		let pre = document.createElement("pre");
		pre.textContent = this.responseText;
		pre.style.overflow = "unset";
		logsContent.appendChild(pre);

		container.appendChild(logsContent);
	}

	let logsReq = new XMLHttpRequest();
	logsReq.addEventListener("load", logsContentHandler);
	logsReq.open("GET", url);
	logsReq.send();
}

function setUserStatus(data) {
	let value = data["value"];
	for (let i in value) {
		if (value[i]["Online"]) {
			document.getElementById("status-" + i).textContent = "Online";
		}
	}
}

var oReq = new XMLHttpRequest();
oReq.addEventListener("load", reqListener);
oReq.open("GET", "/logs/");
oReq.send();
