let lastPath = "";

function getFolders(path) {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;

	xhr.onreadystatechange = function () {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			lastPath = path;
			parseResults(JSON.parse(this.responseText));
		}
	};

	xhr.open("GET", "/api/browse?path=" + path, true);
	xhr.send();
}

function parseResults(parsedResult) {
	let table = document.getElementById("browseBody");
	table.innerHTML = "";
	for (let i = 0; i < parsedResult.length; i++) {
		let tr = document.createElement("tr");
		let td = document.createElement("td");
		let a = document.createElement("a");
		a.href = "#";
		a.className = "d-flex";
		a.innerText = parsedResult[i]["name"];
		a.onclick = function () {
			getFolders(parsedResult[i]["path"]);
		};
		td.appendChild(a);

		tr.appendChild(td);
		table.appendChild(tr);
	}
}

function addPath() {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;

	xhr.onreadystatechange = function () {
		window.location.reload();
	};

	xhr.open("POST", "/api/setPath?path=" + lastPath, true);
	xhr.send();
}
