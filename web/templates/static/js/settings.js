function hideProgress(id) {
    document.getElementById("container-" + id).style.height = "0px";
    document.getElementById("progress-" + id).style.height = "0px";
    document.getElementById("progressbg-" + id).style.height = "0px";
}

function showProgress(id) {
    document.getElementById("container-" + id).style.height = "111px";
    document.getElementById("progress-" + id).style.height = "22px";
    document.getElementById("progressbg-" + id).style.height = "22px";
}

function toggleOnClick(uid) {
    var toggle = document.getElementById("toggle-" + uid);
    if (toggle.classList.contains("active")) {
        hideProgress(uid);
        toggle.classList.remove("active");
        toggle.setAttribute("style", 'transform: rotate(0deg)');
    } else {
        toggle.classList.add("active");
        showProgress(uid);
        toggle.setAttribute("style", 'transform: rotate(180deg)');
    }
}

function bringSettings() {
    var main = document.getElementById('navbar-contents-main');
    var settings = document.getElementById('navbar-contents-settings');
    main.hidden = true;
    settings.hidden = false;
}
bringSettings()


function update_progress(progress, uid) {
    var progressBar = document.getElementById("progress-" + uid);
    var progress_percent = (progress/100)*847;
    var percent = document.getElementById(uid + "-progress-percent");
    percent.innerHTML = progress + "%";
    progressBar.style.width = progress_percent + "px";
}

function stop_task(uid) {
    webStream.send(JSON.stringify({
        'task_id': "stop_task",
        'uid': uid
    }));
}

var start_scan = function(force="false") {
    webStream.send(JSON.stringify({
        'task_id': "start_scan",
        'force': force
    }));
}

var start_regen = function() {
    webStream.send(JSON.stringify({
        'task_id': "start_regen",
    }));
}

function get_running_tasks() {
    webStream.send(JSON.stringify({
        'task_id': "get_running_tasks",
    }));
}

function set_running_tasks(result) {
    Object.keys(result).forEach(function(key) {
        setProgress(result[key]['progress'], key, result[key]['name']);
    });
}


function setProgress(progress, uid, name) {
    if (name == "Scan") {
        update_progress(progress, "scan");
        return;
    }
    if (name == "RegenDB") {
        update_progress(progress, "regen");
        return;
    }
    var root_element = document.getElementById(uid);
    if (root_element == null) {
        addProgressBar(uid, name);
    } else {
        update_progress(progress, uid);
    }
}

function addProgressBar(uid, name){
    var tasks = document.getElementById('tasks')
    var template = document.createElement('template');
    var htmlText = ' <div id=' + uid +' class="relative-task"> <div class="d-inline-flex"> <span>' + name + '</span> <span id="' + uid +'-progress-percent" class="text-right bold">0%</span> <span class="text-right">Progress details</span> <i id="toggle-' + uid + '" class="fa fa-angle-down text-right dropdown-icon" onclick="toggleOnClick(\'' + uid +'\')"></i></div> <div> <div id="container-' + uid + '" class="drop-container"> <div class="task-progress-container"> <div id="progressbg-' + uid + '" class="progress-bar-bg"></div> <div id="progress-' + uid + '" class="progress-bar"></div> </div> <div class="d-flex progress-footer"><span>Time started</span><span class="footer-right">Status</span></div> </div> </div> </div>'
    template.innerHTML = htmlText.trim();
    var clone2 = document.importNode(template.content, true);
    tasks.append(clone2)
}

function filepath_poll(path)  {
        webStream.send(JSON.stringify({
            'task_id': "filepath_poll",
            'path': path
        }));
}

const webStream = new WebSocket(
        'ws://'
        + window.location.host
        + '/ws/Jizzberry/settings'
        + '/'
    );

webStream.onmessage = function(e) {
    const data = JSON.parse(e.data);

    if (data.task_id == "filepath_poll") {
        console.log(data.result);
    } else if (data.task_id == "save_filepath") {
        add_filepath(data.result);
    } else if (data.task_id == "progress_update") {
        setProgress(data.progress, data.uid, data.name)
    } else if (data.task_id == "get_running_tasks") {
        set_running_tasks(data.result);
    }
};

webStream.onclose = function(e) {
    console.error('Chat socket closed unexpectedly');
};

webStream.onopen = function() {
    filepath_poll("")
    get_running_tasks();
}

