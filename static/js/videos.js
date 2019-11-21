function downloadVideo() {
    startSpinner("download-video-spinner")
    let videoURL = document.getElementById("video-url").value
    let downloadMode = document.getElementById("download-mode").value
    let fileExtension = document.getElementById("file-extension").value
    let downloadQuality = document.getElementById("download-quality").value
    let downloadPath = document.getElementById("download-path").value

    console.log(downloadPath)

    let videoData = {
        videoURL,
        downloadMode,
        fileExtension,
        downloadQuality,
        downloadPath,
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
          <p>Download Path: ${video.DownloadPath}</p>
        </div>
      </div>
    </div>`
    })
}
