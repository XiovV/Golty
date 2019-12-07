function downloadVideo() {
    startSpinner("download-video-spinner");
    let videoURL = document.getElementById("video-url").value;
    let downloadMode = document.getElementById("download-mode").value;
    let fileExtension = document.getElementById("file-extension").value;
    let downloadQuality = document.getElementById("download-quality").value;
    let downloadPath = document.getElementById("output-path-indicator").innerText;

    console.log(downloadPath);

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
            handleResponse(res);
            stopSpinner("download-video-spinner");
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

function displayVideos(videos) {
    document.getElementById("accordion").innerHTML = "";
    console.log(videos);
    videos.forEach((video, index) => {
        console.log(video);
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

function customYtdl(checkboxId) {
    document.getElementById("download-path").disabled = false
    if (checkboxId === "custom-download-output") {
        document.getElementById("download-path").placeholder = "default: /videos/"
    } else if (checkboxId === "custom-ytdl-output") {
        document.getElementById("download-path").placeholder = "default: /videos/%(uploader)s/audio/%(title)s.%(ext)s"
    }
}

function changeOutputPathIndicator(id) {
    document.getElementById("output-path-indicator").innerHTML = "";
    let downloadPathRadio = document.getElementById("custom-download-output").checked;
    let youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
    let downloadMode = document.getElementById("download-mode").value;
    let input = document.getElementById(id).value;
    if (downloadMode === "Audio Only") {
        if (downloadPathRadio === true) {
            document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/audio/%(title)s.%(ext)s"
        } else if (youtubedlOutputRadio === true) {
            document.getElementById("output-path-indicator").innerHTML = input
        }
    } else if (downloadMode === "Video And Audio") {
        if (downloadPathRadio === true) {
            document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/video/%(title)s.%(ext)s"
        } else if (youtubedlOutputRadio === true) {
            document.getElementById("output-path-indicator").innerHTML = input
        }
    }
}