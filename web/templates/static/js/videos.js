let element = document.getElementById("videoPlayer")
let player;
let playable = false;

const sessionsSocket = new WebSocket(
    "ws://" + window.location.host + "/ws/session" + ""
);

sessionsSocket.onopen = function () {
    let xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (element.canPlayType(this.responseText) === "maybe" || element.canPlayType(this.responseText) === "probably") {
                playable = true
            }
            sessionsSocket.send(JSON.stringify({
                type: "getStreamURL",
                data: JSON.stringify({
                    scene_id: sceneId,
                    playable: playable,
                    start_time: 0
                })
            }))
        }
    }
    xhr.open("GET", '/api/getMimeType?scene_id=' + sceneId, true);
    xhr.send();
}

sessionsSocket.onmessage = function (e) {
    let parsed = JSON.parse(e.data);

    switch (parsed.type) {
        case "getStreamURL":
            let data = JSON.parse(parsed.data)
            if (initial) {
                let source = document.createElement('source');
                source.setAttribute('src', data['URL'])
                source.setAttribute('type', data['MimeType']);
                element.appendChild(source)

                player = new Plyr(element, {duration: duration});

                if (!playable) {
                    player.on('ready', event => {
                        player.on("seeking", function () {
                            console.log("called seek")
                            sessionsSocket.send(JSON.stringify({
                                type: "getStreamURL",
                                data: JSON.stringify({
                                    scene_id: sceneId,
                                    playable: false,
                                    start_time: player.currentTime
                                })
                            }))
                            offset = player.currentTime
                        })
                    });
                    initial = false
                } else {
                    player.source = {
                        type: "video",
                        sources: [
                            {
                                src: data['URL'],
                                type: data['MimeType']
                            },
                        ]
                    }
                    player.offset = offset;
                }
            }
    }

    console.log(parsed)
};

sessionsSocket.onclose = function (e) {
    console.error(e);
};

$("body").on("click", "#parser-toggle", function () {
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
    let title = document.getElementById("metadata-title");
    let url = document.getElementById("metadata-url");
    let date = document.getElementById("metadata-date");

    postMetadata(title.value, url.value, date.value, actors, tags, studios)
}

function postMetadata(title, url, date, actors, tags, studios) {
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
        title: title,
        url: url,
        date: date,
        actors: actors,
        tags: tags,
        studios: studios,
    }));
}

function getQueryResults(term) {
    let container = document.getElementById("query-body")
    container.innerHTML = ""
    if (term === undefined || term === "") {
        let title = document.getElementById("video-title");
        term = title.textContent;
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

        function submitResult(r) {
            postMetadata(r["name"], r["url"], "", r["actors"], r["tags"], [])
        }

        let li = document.createElement("li")
        container.appendChild(li)

        {
            let div = document.createElement("div")
            li.appendChild(div)

            let title = document.createElement("p")
            title.className = "m-0 text-muted"
            title.innerText = result[i]["name"]
            div.appendChild(title)

            let websiteButton = document.createElement("button")
            websiteButton.className = "badge badge-success p-1 ml-3 btn text-white"
            websiteButton.style.boxShadow = "none"
            websiteButton.style.outline = "none"
            websiteButton.innerText = result[i]["website"]
            title.appendChild(websiteButton)
        }

        {
            let submitDiv = document.createElement("div")
            submitDiv.className = "d-flex justify-content-end"
            li.appendChild(submitDiv)

            let submitBtton = document.createElement("button")
            submitBtton.className = "badge badge-success p-1 btn text-white"
            submitBtton.style.boxShadow = "none"
            submitBtton.style.outline = "none"
            submitBtton.onclick = function () {
                submitResult(result[i])
            }
            submitDiv.appendChild(submitBtton)

            let submitIcon = document.createElement("i")
            submitIcon.className = "fa fa-check"
            submitIcon.style.fontSize = "18px"
            submitIcon.style.width = "2rem"
            submitBtton.appendChild(submitIcon)
        }

        {
            let linkText = document.createElement("p")
            linkText.style.fontSize = "14px"
            linkText.className = "h6 m-0 text-muted"
            linkText.innerText = "Path:"
            li.appendChild(linkText)

            let link = document.createElement("a")
            link.className = "text-muted query-link"
            link.href = result[i]["url"]
            link.innerText = result[i]["url"]
            li.appendChild(link)

        }

        let ul = document.createElement("ul")
        ul.className = "list-group p-0"
        container.appendChild(ul)

        let actorsDiv = document.createElement("div")
        ul.appendChild(actorsDiv)

        let tagsDiv = document.createElement("div")
        ul.appendChild(tagsDiv)

        function addSearch(className, container, addFunc, placeholder) {
            let searchLi = document.createElement("li")
            searchLi.className = "list-group-item d-flex justify-content-between p-0 border-0"
            let searchDiv = document.createElement("div")
            searchDiv.className = "input-group mt-3 p-0"
            searchLi.appendChild(searchDiv)

            let searchInput = document.createElement("input")
            searchInput.type = "search"
            searchInput.className = className + " form-control no-border"
            searchInput.placeholder = placeholder
            searchDiv.appendChild(searchInput)

            let addButtonDiv = document.createElement("div")
            addButtonDiv.className = "input-group-append"
            addButtonDiv.style.width = "2.5rem"
            searchDiv.appendChild(addButtonDiv)

            let addButton = document.createElement("button")
            addButton.style.borderRadius = "0 3px 3px 0"
            addButton.className = "p-0 input-group-text btn"
            addButton.onclick = function () {
                addFunc(searchInput.value.toString());
                searchInput.value = "";
                $('.typeaheadActors').typeahead('val', '');
            }
            addButtonDiv.appendChild(addButton)

            let addButtonSpan = document.createElement("span")
            addButtonSpan.style.boxShadow = "none"
            addButtonSpan.style.outline = "none"
            addButtonSpan.style.paddingTop = "10px"
            addButtonSpan.style.paddingBottom = "10px"
            addButtonSpan.style.borderRadius = "0 3px 3px 0"
            addButtonSpan.className = "input-group-text btn"
            addButton.appendChild(addButtonSpan)

            let addButtonIcon = document.createElement("i")
            addButtonIcon.className = "fas fa-plus fa-fw"
            addButtonSpan.appendChild(addButtonIcon)

            container.appendChild(searchLi)
        }

        function addData(className, array, removeFunc, addFunc, placeholder, container) {
            container.innerHTML = ""
            addSearch(className, container, addFunc, placeholder)
            if (array !== undefined && array != null) {
                for (let j = 0; j < array.length; j++) {
                    let actors = document.createElement("li")
                    actors.className = "list-group-item d-flex justify-content-between p-0 border-0"

                    let div = document.createElement("div")
                    div.className = "input-group p-0"

                    actors.appendChild(div)

                    let field = document.createElement("input")
                    field.className = "form-control bg-white no-border"
                    field.type = "text"
                    field.disabled = true
                    field.value = array[j]

                    div.appendChild(field)

                    let buttonDiv = document.createElement("div")
                    buttonDiv.className = "input-group-append"
                    buttonDiv.style.width = "2.5rem"

                    let button = document.createElement("button")
                    button.className = "p-0 input-group-text btn";
                    button.style.borderRadius = "0px 3px 3px 0px"
                    button.onclick = function () {
                        removeFunc(j)
                    }

                    let span = document.createElement("span")
                    span.className = "input-group-text btn"
                    span.style.paddingTop = "10px"
                    span.style.paddingBottom = "10px"
                    span.style.borderRadius = "0px 3px 3px 0px"

                    let icon = document.createElement("i")
                    icon.className = "fas fa-times fa-fw"

                    span.appendChild(icon)

                    buttonDiv.appendChild(button)
                    button.appendChild(span)
                    div.appendChild(buttonDiv)

                    container.appendChild(actors)
                }
            }
        }

        function refreashAll() {
            function refreshActors() {
                function remove(j) {
                    result[i]["actors"].splice(j, 1)
                    refreashAll()
                }

                function add(value) {
                    if (result[i]["actors"] == null) {
                        result[i]["actors"] = [value]
                    } else {
                        result[i]["actors"].push(value)
                    }
                    refreashAll()
                }

                addData("typeaheadActors", result[i]["actors"], remove, add, "Search Actors...", actorsDiv)
            }

            function refreshTags() {
                function remove(j) {
                    result[i]["tags"].splice(j, 1)
                    refreashAll()
                }

                function add(value) {
                    if (result[i]["tags"] == null) {
                        result[i]["tags"] = value
                    } else {
                        result[i]["tags"].push(value)
                    }
                    refreashAll()
                }

                addData("typeaheadTags", result[i]["tags"], remove, add, "Search Tags...", tagsDiv)
            }

            refreshActors()
            refreshTags()

            $('.typeaheadActors').typeahead(
                {
                    hint: true,
                    minLength: 0,
                    highlight: true,
                },
                {
                    name: 'items',
                    displayKey: 'name',
                    source: actorsData.ttAdapter(),
                    templates: {
                        suggestion: Handlebars.compile("<p><button style='display: block;'>{{name}}</button></p>")
                    },
                    limit: 5
                },
            );

            $('.typeaheadTags').typeahead(
                {
                    hint: true,
                    minLength: 0,
                    highlight: true,
                },
                {
                    name: 'items',
                    displayKey: 'name',
                    source: tagsData.ttAdapter(),
                    templates: {
                        suggestion: Handlebars.compile("<p><button style='display: block;'>{{name}}</button></p>")
                    },
                    limit: 5
                },
            );
        }

        refreashAll()
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

$('#query-search-button').on('click', function () {
    getQueryResults($('#query-search').val())
})

