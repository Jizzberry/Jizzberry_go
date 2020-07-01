$("body").on("click","#parser-toggle", function(){
    var toggle = document.getElementById('parser-toggle');
    var form = document.getElementById('metadata-parser')
    var viewer = document.getElementById("metadata-viewer");
    if (toggle.classList.contains("active")) {
        toggle.classList.remove("active");
        form.className = "d-none";
        viewer.hidden = false;
    } else {
        toggle.classList.add("active")
        form.className = "d-flex flex-column mb-3"
        viewer.hidden = true;
    }
});


let actors = [];
let tags = [];
let studios = [];

let urlSplit = window.location.href.split("/");
let sceneId = urlSplit[urlSplit.length - 1]

getArrays()


const actorsData = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
        url: "/api/actors?name=%QUERY",
        wildcard: '%QUERY',
        filter: function (response) {
            console.log(response);
            return response;
        }
    }
});

const tagsData = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
        url: "/api/tags?name=%QUERY",
        wildcard: '%QUERY',
        filter: function (response) {
            console.log(response);
            return response;
        }
    }
});

const studiosData = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
        url: "/api/studios?name=%QUERY",
        wildcard: '%QUERY',
        filter: function (response) {
            console.log(response);
            return response;
        }
    }
});

studiosData.initialize();
tagsData.initialize();
actorsData.initialize();


function getArrays() {
    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            let response = JSON.parse(this.responseText);
            for (const [key, value] of Object.entries(response)) {
                actors = value['actors'].toString().split(", ")
                tags = value['tags'].toString().split(", ")
                studios = value['studios'].toString().split(", ")
            }

            function removeEmpty(array) {
                return array.filter(function (element) {
                    return (element !== "" && element !== null)
                })
            }
            actors = removeEmpty(actors);
            studios = removeEmpty(studios);
            tags = removeEmpty(tags);
        }
    }

    xhr.open("GET", '/api/files?generated_id='+sceneId, true);
    xhr.send();
}

function refreshModalData(arrayName) {
    let array = getArray(arrayName)
    let inputField = document.getElementById("multiselector-search")
    let addButton = document.getElementById("multiselector-add")
    addButton.onclick = function () {
        addDataArray(arrayName, inputField.value)
    }
    let ul = document.getElementById("multiselector-container")

    if (ul != null) {
        ul.innerHTML = "";
        for (let i = 0; i < array.length; i++) {
            let li = document.createElement("li")
            li.style.fontSize = "15px"
            li.className = "list-group-item d-flex justify-content-between p-2 px-2"

            let span = document.createElement("span")
            span.style.paddingTop = ".15rem"
            span.textContent = array[i]

            function removeElement() {
                array.splice(i, 1)
                refreshModalData(arrayName)
            }

            let button = document.createElement("button")
            button.style.verticalAlign = "middle"
            button.type = "button"
            button.className = "close"
            button.onclick = removeElement

            let buttonSpan = document.createElement("span")
            buttonSpan.className = "align-text-top"
            buttonSpan.innerHTML = "&times;"

            ul.appendChild(li)
            li.appendChild(span)
            li.appendChild(button)
            button.appendChild(buttonSpan)
        }
    }
}

function getArray(arrayName) {
    switch (arrayName) {
        case "actors":
            return actors
        case "tags":
            return tags;
        case "studios":
            return studios;
        default:
            return null;
    }
}

function addDataArray(arrayName, element) {
    switch (arrayName) {
        case "actors":
            actors.push(element)
            refreshModalData("actors")
            return;
        case "tags":
            tags.push(element);
            refreshModalData("tags")
            return;
        case "studios":
            studios.push(element)
            refreshModalData("studios")
            return;
        default:
            return;
    }
}

function openActorsModal() {
    refreshModalData("actors");
    reinitializeTypeahead("actors")
}

function openTagsModal() {
    refreshModalData("tags");
    reinitializeTypeahead("tags")
}

function openStudiosModal() {
    refreshModalData("studios");
    reinitializeTypeahead("studios")
}

function saveMetadata() {
    var title = document.getElementById("metadata-title");
    var url = document.getElementById("metadata-url");
    var date = document.getElementById("metadata-date");

    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            window.location.reload()
        }
    }

    xhr.open("POST", '/api/metadata', true);
    xhr.send(JSON.stringify({
        generated_id: sceneId,
        title: title.value,
        url: url.value,
        date: date.value,
        actors: actors,
        tags: tags,
        studios: studios,
    }));
}

function getQueryResults(term) {
    let container = document.getElementById("query-body")
    container.innerHTML = ""
    if (term === undefined) {
        let title = document.getElementById("video-title");
        term = title.textContent;
        console.log(term)
    }
    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            refreshQueryModal(JSON.parse(this.responseText))
        }
    }

    xhr.open("GET", '/api/queryScrapers?term=' + encodeURIComponent(term.trim()), true);
    xhr.send();

}

function refreshQueryModal(result) {
    console.log(result)
    let container = document.getElementById("query-body")
    container.innerHTML = ""
    for (let i = 0; i < result.length; i++) {
        let li = document.createElement("li")
        li.innerHTML = '<div><p class="m-0 text-muted">' + result[i]["Name"] + '<button style="box-shadow: none; outline: none;"class="badge badge-success p-1 ml-3 btn text-white">' + result[i]["Website"] + '</button></p></div><div class="d-flex justify-content-end"><button style="box-shadow: none; outline: none;" class="badge badge-success p-1 btn text-white"><i style="font-size: 18px; width: 2rem;" class="fa fa-check" aria-hidden="true"></i></button></div><p style="font-size: 14px;" class="h6 m-0 text-muted">Path:</p><a class="text-muted" style="outline: none;color: black;text-decoration: none;border: none;box-shadow: none;font-size: 14px;" href="#">' + result[i]["Url"] + '</a>';
        container.appendChild(li)

        let ul = document.createElement("ul")
        ul.className = "list-group p-0"

        let searchLi = document.createElement("li")
        let searchDiv = document.createElement("div")
        searchDiv.innerHTML = '<div class="input-group mt-3 p-0"><input style="box-shadow: none; outline: none;" type="search" class="form-control no-border" placeholder="Enter..." /><div class="input-group-append" style="width: 2.5rem;"><button style="border-radius: 0 3px 3px 0;" class="p-0 input-group-text btn"><span style="box-shadow: none; outline: none; padding-top: 10px;padding-bottom: 10px; border-radius: 0 3px 3px 0;" class="input-group-text btn"><i class="fas fa-plus fa-fw"></i></span></button></div></div>';

        searchLi.appendChild(searchDiv)

        ul.appendChild(searchLi)
        container.appendChild(ul)

        for (let j = 0; j < result[i]["Actors"].length; j++) {
            let actors = document.createElement("li")
            actors.className = "list-group-item d-flex justify-content-between p-0 border-0"

            let div = document.createElement("div")
            div.className = "input-group p-0"

            actors.appendChild(div)

            let field = document.createElement("input")
            field.className = "form-control bg-white no-border"
            field.type = "text"
            field.disabled = true
            field.value = result[i]["Actors"][j]

            div.appendChild(field)

            let buttonDiv = document.createElement("div")
            buttonDiv.className = "input-group-append"
            buttonDiv.style.width = "2.5rem"

            let button = document.createElement("button")
            button.className = "p-0 input-group-text btn";
            button.style.borderRadius = "0px 3px 3px 0px"

            let span = document.createElement("span")
            span.className = "input-group-text btn"
            span.style.paddingTop = "10px"
            span.style.paddingBottom = "10px"
            span.style.borderRadius = "0px 3px 3px 0px"

            let icon = document.createElement("i")
            icon.className = "fas fa-times fa-fw"

            span.appendChild(icon)

            buttonDiv.appendChild(button)
            buttonDiv.appendChild(span)
            div.appendChild(buttonDiv)

            ul.appendChild(actors)
        }
    }
}

function reinitializeTypeahead(type) {
    let source
    switch (type) {
        case "actors":
            source = actorsData.ttAdapter();
            break;
        case "tags":
            source = tagsData.ttAdapter();
            break
        case "studios":
            source = studiosData.ttAdapter();
            break
        default:
            return
    }
    $('#multiselector-search').typeahead("destroy");

    $('#multiselector-search').typeahead(
        {
            hint: true,
            minLength: 0,
            highlight: true,
        },
        {
            name: 'items',
            displayKey: 'name',
            source: source,
            templates: {
                suggestion: Handlebars.compile("<p><button style='display: block;'>{{name}}</button></p>")
            },
            limit: 5
        },
    );
}

