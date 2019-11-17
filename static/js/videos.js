function downloadVideo() {
    startSpinner("download-video-spinner")
    let videoURL = document.getElementById("video-url").value
    let downloadMode = document.getElementById("download-mode").value
    let fileExtension = document.getElementById("file-extension").value
    let downloadQuality = document.getElementById("download-quality").value


    let videoData = {
        videoURL,
        downloadMode,
        fileExtension,
        downloadQuality,
    };

    const options = {
        method: "POST",
        body: JSON.stringify(videoData),
        headers: new Headers({
            "Content-Type": "application/json"
        })
    };

    fetch("/api/download-video", options)
        .then(res => res.json())
        .then(res => {
            handleResponse(res)
            stopSpinner("download-video-spinner")
            getVideos()
        });
}

function getVideos() {
    fetch("/api/get-videos")
        .then(res => res.json())
        .then(videos => {
            displayVideos(videos);
            getVersion();
        });
}

function displayErrorMessage(message) {
    let error = document.getElementById("error");
    error.innerHTML = ""
    error.classList.remove("d-none");
    error.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displaySuccessMessage(message) {
    let success = document.getElementById("success");

    if (success) {
        console.log("DISPLAY SUCCESS ALERT")

        success.innerHTML = ""
        success.classList.remove("d-none");
        success.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
    } else {
        console.log("CREATE SUCCESS ALERT")
        let alertsDiv = document.getElementById("alerts").innerHTML
        alertsDiv += `<div class="alert alert-success alert-dismissible mt-3" id="success" role="alert">
                    ${message}
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                    </button>
                  </div>`
    }
}

function displayWarningMessage(message) {
    let warning = document.getElementById("warning");
    warning.innerHTML = ""
    warning.classList.remove("d-none");
    warning.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displayVideos(videos) {
    document.getElementById("accordion").innerHTML = ""
    console.log(videos)
    videos.forEach((video, index) => {
        console.log(video)
        document.getElementById("accordion").innerHTML += `<div class="mb-2 p-2 card">
      <h5 class="mb-0">
        <button class="btn btn-link dropdown-toggle" data-toggle="collapse" data-target="#collapse${index}" aria-expanded="true" aria-controls="collapse${index}" id=${video.ChannelURL}listElem>
          ${video.VideoURL}
        </button>
      </h5>
  
      <div id="collapse${index}" class="collapse" aria-labelledby="heading${index}" data-parent="#accordion">
        <div class="panel-body ml-2">
          <p>Downloaded As: ${video.DownloadMode}</p>
          <p>File Extension: ${video.FileExtension}</p>
        </div>
      </div>
    </div>`
    })
}
function handleResponse(res) {
    if (res.Type == "Success") {
        displaySuccessMessage(res.Message)
    } else if (res.Type == "Error") {
        displayErrorMessage(res.Message)
    } else if (res.Type == "Warning") {
        displayWarningMessage(res.Message)
    }
}

function startSpinner(id) {
    spinner = document.getElementById(id);
    spinner.classList.remove("d-none");
}

function stopSpinner(id) {
    spinner = document.getElementById(id);
    spinner.classList.add("d-none")
}

function changeExtension() {
    console.log("change ext")
    let downloadMode = document.getElementById("download-mode").value
    let fileExtensions = document.getElementById("file-extension")
    let downloadQualities = document.getElementById("download-quality")
    if (downloadMode == "Audio Only") {
        fileExtensions.options[0].value = "m4a"
        fileExtensions.options[0].text = "m4a"
        fileExtensions.options[1].value = "mp3"
        fileExtensions.options[1].text = "mp3"
        downloadQualities.options[0].value = "best"
        downloadQualities.options[0].text = "best"
        downloadQualities.options[1].value = "medium"
        downloadQualities.options[1].text = "medium"
        downloadQualities.options[2].value = "worst"
        downloadQualities.options[2].text = "worst"
    } else if (downloadMode == "Video And Audio") {
        fileExtensions.options[0].value = "any"
        fileExtensions.options[0].text = "any (recommended for now)"
        fileExtensions.options[1].value = "mp4"
        fileExtensions.options[1].text = "mp4"
        // fileExtensions.options[2].value = ".mkv"
        // fileExtensions.options[2].text = ".mkv"

        downloadQualities.options[0].value = "best"
        downloadQualities.options[0].text = "best"
        downloadQualities.options[1] = null
        downloadQualities.options[2].value = "worst"
        downloadQualities.options[2].text = "worst"
    }
}