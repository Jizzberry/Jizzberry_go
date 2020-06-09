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

let actors = ["test", "hello", "bye"];
let tags = ["test", "hello", "bye"];
let studios = ["test", "hello", "bye"];

let urlSplit = window.location.href.split("/");
let sceneId = urlSplit[urlSplit.length - 1]

getArrays()

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

function openModal(array) {
    let body = document.getElementById("multiselector-modal-body")

    if (body != null) {
        body.innerHTML = "";
        console.log(array)
        for (let i = 0; i < array.length; i++) {
            let text = document.createElement("span");
            let remove = document.createElement("button")
            let div = document.createElement("div");
            text.textContent = array[i];
            remove.onclick = removeElement;

            function removeElement() {
                array.splice(i, 1)
                console.log(array)
                openModal(array)
            }

            // div.appendChild(text)
            // div.appendChild(remove)
            // body.appendChild(div)
        }
    }
}

function openActorsModal() {
    openModal(actors);
}

function openTagsModal() {
    openModal(tags);
}

function openStudiosModal() {
    openModal(studios);
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