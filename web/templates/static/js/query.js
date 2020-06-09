function fetchQueryResults(term) {
    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            console.log(this.responseText)
            parseQueryResults(JSON.parse(this.responseText))
        }
    }

    xhr.open("GET", '/api/queryScrapers?term=' + term, true);
    xhr.send();
}

function parseQueryResults(resultSet) {
    let container = document.getElementById("queryModalBody");

    let ul = document.createElement("ul")
    ul.className = "list-group"

    for (let i = 0; i < resultSet.length; i++) {
        let li = document.createElement("li")
        li.className = "list-group-item"

        let heading = document.createElement("h6")
        heading.textContent = resultSet[i]['Name']

        let url = document.createElement("span")
        url.style.fontSize = "12px"
        url.innerText = resultSet[i]['Url']

        let tableContainer = document.createElement("div")
        tableContainer.className ="table-responsive"

        let actorTable = document.createElement("table")
        actorTable.className = "table"

        let tableBody = document.createElement("tbody")


        for (let j = 0; j < resultSet[i]['Actors'].length; j++) {
            let tr = document.createElement("tr")
            tr.innerText = resultSet[i]['Actors'][j]
            tableBody.appendChild(tr)
        }

        actorTable.appendChild(tableBody)
        tableContainer.appendChild(tableBody)

        li.appendChild(heading)
        li.appendChild(tableContainer)

        ul.appendChild(li)
    }
    container.appendChild(ul)
}
