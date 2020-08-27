function fetchQueryResults(term) {
	let xhr = new XMLHttpRequest();
	xhr.withCredentials = true;

	xhr.onreadystatechange = function () {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			console.log(this.responseText);
			parseQueryResults(JSON.parse(this.responseText));
		}
	};

	xhr.open("GET", "/api/queryScrapers?term=" + term, true);
	xhr.send();
}

function parseQueryResults(resultSet) {
	let container = document.getElementById("queryModalBody");

	let ul = document.createElement("ul");
	ul.className = "list-group";

	for (let i = 0; i < resultSet.length; i++) {

	}
	container.appendChild(ul);
}
