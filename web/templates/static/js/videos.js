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

function refreshModalData(array) {
    let body = document.getElementById("multiselector-form")

    if (body != null) {
        body.innerHTML = "";
        for (let i = 0; i < array.length; i++) {
            let holder = document.createElement("div");

            let inputGrp = document.createElement("div");
            inputGrp.className = "input-group";

            let inputField = document.createElement("input");
            inputField.type = "text";
            inputField.className = "form-control"
            inputField.disabled = true;
            inputField.placeholder = array[i];

            inputGrp.appendChild(inputField)

            let buttonHolder = document.createElement("div")
            buttonHolder.className = "input-group-prepend button-holder"

            function removeElement() {
                array.splice(i, 1)
                console.log(array)
                refreshModalData(array)
            }

            let button = document.createElement("button")
            button.className = "input-group-text bg-light";
            button.onclick = removeElement;

            buttonHolder.appendChild(button)

            holder.appendChild(inputGrp)
            holder.appendChild(buttonHolder)

            body.appendChild(holder)
        }
    }
}

function addDataArray(arrayName, element) {
    switch (arrayName) {
        case "actors":
            actors.push(element)
            refreshModalData(actors)
            return;
        case "tags":
            tags.push(element);
            refreshModalData(tags)
            return;
        case "studios":
            studios.push(element)
            refreshModalData(studios)
            return;
        default:
            return;
    }
}

function openActorsModal() {
    refreshModalData(actors);
    reinitializeTypeahead("actors")
}

function openTagsModal() {
    refreshModalData(tags);
    reinitializeTypeahead("tags")
}

function openStudiosModal() {
    refreshModalData(studios);
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
    $('#multiselector-modal-body .typeahead').typeahead("destroy");

    $('#multiselector-modal-body .typeahead').typeahead(
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
                suggestion: Handlebars.compile("<p><button onclick='addDataArray(\"" + type + "\", \"{{name}}\")' style='display: block;'>{{name}}</button></p>")
            },
            limit: 8
        },
    );
}

