
function updateCheckingInterval() {
  startSpinner("update-checking-interval-spinner")
  let checkingInterval
  let checkingIntervalInput = document.getElementById("checking-interval").value
  let time = document.getElementById("time").value

  if (time == "minutes") {
    checkingInterval = checkingIntervalInput 
  } else if (time == "hours") {
    checkingInterval = checkingIntervalInput * 60
  } else if (time == "days") {
    checkingInterval = checkingIntervalInput * 1440
  }

  interval = {
    checkingInterval
  }

  const options = {
    method: "POST",
    body: JSON.stringify(interval),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/update-checking-interval-playlists", options)
    .then(res => res.json())
    .then(res => {
      handleResponse(res)
      stopSpinner("update-checking-interval-spinner")
    });
}

function addPlaylist() {
  startSpinner("add-playlist-spinner")
  let downloadEntire = document.querySelector('#download-entire-playlist').checked;
  let URL = document.getElementById("playlist-url").value;
  let downloadMode = document.getElementById("download-mode").value;
  let fileExtension = document.getElementById("file-extension").value;
  let downloadQuality = document.getElementById("download-quality").value;
  let downloadPath = document.getElementById("output-path-indicator").innerText;

  let type = "Playlist";

  let playlistData = {
    URL,
    downloadMode,
    fileExtension,
    downloadQuality,
    downloadEntire,
    downloadPath,
    type,
  };

  const options = {
    method: "POST",
    body: JSON.stringify(playlistData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/add", options)
    .then(res => res.json())
    .then(res => {
      handleResponse(res)
      stopSpinner("add-playlist-spinner")
      getPlaylists()
    });
}

function checkAll() {
  startSpinner("check-all-spinner")
  let channelData = {
    Type: "playlists"
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/check-all", options)
      .then(res => res.json())
      .then(res => {
        handleResponse(res)
        stopSpinner("check-all-spinner")
        getPlaylists()
      });
}

function getPlaylists() {
  let channelData = {
    Type: "playlists"
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/get", options)
      .then(res => res.json())
      .then(playlists => {
        displayPlaylists(playlists);
        getVersion();
      });
}

function checkPlaylist(id) {
  startSpinner(id+"-spinner");
  let URL = id;
  let downloadMode = document.getElementById("download-mode").value;
  let fileExtension = document.getElementById("file-extension").value;
  let downloadQuality = document.getElementById("download-quality").value;

  let channelData = {
    URL,
    downloadMode,
    fileExtension,
    downloadQuality
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/check-playlist", options)
    .then(res => res.json())
    .then(res => {
      console.log(res);
      stopSpinner(id+"-spinner")
      if (res.Type == "Success") {
        if (res.Key == "NO_NEW_VIDEOS") {
          displayWarningMessage(res.Message);
          getPlaylists()
        } else if (res.Key == "NEW_VIDEO_DETECTED") {
          displaySuccessMessage(res.Message);
          getPlaylists()
        }
      } else if (res.Type == "Error") {
        if (res.Key == "ERROR_DOWNLOADING_VIDEO") {
          displayErrorMessage(res.Message);
        }
      }
    });
}

function deletePlaylist(id) {
  let playlistURL = {
    URL: id,
    Type: "Playlist"
  };

  const options = {
    method: "POST",
    body: JSON.stringify(playlistURL),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/delete", options)
    .then(res => res.json())
    .then(res => {
      handleResponse(res);
      getPlaylists()
    });
}

function displayErrorMessage(message) {
  let error = document.getElementById("error");
  error.innerHTML = "";
  error.classList.remove("d-none");
  error.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displaySuccessMessage(message) {
  let success = document.getElementById("success");

  if (success) {
    console.log("DISPLAY SUCCESS ALERT");

    success.innerHTML = "";
    success.classList.remove("d-none");
    success.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
  } else {
    let alertsDiv = document.getElementById("alerts").innerHTML;
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
  warning.innerHTML = "";
  warning.classList.remove("d-none");
  warning.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displayPlaylists(playlists) {
  document.getElementById("accordion").innerHTML = "";
  playlists.forEach((playlist, index) => {
    console.log(playlist)
    document.getElementById("accordion").innerHTML += `<div class="mb-2 p-2 card">
      <h5 class="mb-0">
        <button class="btn btn-link dropdown-toggle" data-toggle="collapse" data-target="#collapse${index}" aria-expanded="true" aria-controls="collapse${index}" id=${playlist.URL}listElem>
          ${playlist.Name}
        </button><button class="btn btn-danger float-right ml-2" id="${playlist.URL +
        "delTarget"}" onClick="deletePlaylist(this.id)">&times</button><button class="btn btn-primary float-right" id="${playlist.URL}" onClick="checkPlaylist(this.id)">Check<div id="${playlist.URL}-spinner" class="spinner-border align-middle ml-2 d-none"></div></button>
      </h5>
  
      <div id="collapse${index}" class="collapse" aria-labelledby="heading${index}" data-parent="#accordion">
        <div class="panel-body ml-2">
          Latest Download: <a href=https://www.youtube.com/watch?v=${playlist.LatestDownloaded} target="_blank">https://www.youtube.com/watch?v=${playlist.LatestDownloaded}</a>
          <p>Download Mode: ${playlist.DownloadMode}</p>
          <p>Last Checked: ${playlist.LastChecked}</p>
          <p>Preferred Extension For Audio: ${playlist.PreferredExtensionForAudio}
          <p>Preferred Extension For Video: ${playlist.PreferredExtensionForVideo}
          <p>Download Path: ${playlist.DownloadPath}</p>
          <br>
          <button class="btn btn-link dropdown-toggle" type="button" data-toggle="collapse" data-target="#history${index}" aria-expanded="false" aria-controls="history${index}">
            Download History
          </button>
          <div class="collapse" id="history${index}">
            <div class="card card-body" id="dlhistory${playlist.Name}">
            </div>
          </div>
        </h5>
        </div>
      </div>
    </div>`
    displayDownloadHistory(playlist.Name, playlist.DownloadHistory)
  })
}

function displayDownloadHistory(channelName, downloadHistory) {
  let historyBox = document.getElementById("dlhistory"+channelName)
  console.log(historyBox)
  downloadHistory.forEach(video => {
    historyBox.innerHTML += `<br> <a href=https://www.youtube.com/watch?v=${video} target="_blank">https://www.youtube.com/watch?v=${video}</a>` 
  })
}

function changeExtension() {
  let downloadMode = document.getElementById("download-mode").value;
  let fileExtensions = document.getElementById("file-extension");
  let downloadQualities = document.getElementById("download-quality");
  let input = document.getElementById("download-path").value;
  if (downloadMode == "Audio Only") {
    document.getElementById("download-path").placeholder = "default: /playlists/%(uploader)s/audio/%(title)s.%(ext)s";
    if (input.length > 0) {
      downloadPathRadio = document.getElementById("custom-download-output").checked;
      youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
      if (downloadPathRadio == true) {
        document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/audio/%(title)s.%(ext)s"
      } else if (youtubedlOutputRadio == true) {
        document.getElementById("output-path-indicator").innerHTML = input
      }
    } else {
      document.getElementById("output-path-indicator").innerHTML = "/playlists/%(uploader)s/audio/%(title)s.%(ext)s"
    }

    fileExtensions.options[0].value = "m4a";
    fileExtensions.options[0].text = "m4a";
    fileExtensions.options[1].value = "mp3";
    fileExtensions.options[1].text = "mp3";
    downloadQualities.options[0].value = "best";
    downloadQualities.options[0].text = "best";
    downloadQualities.options[1].value = "medium";
    downloadQualities.options[1].text = "medium";
    downloadQualities.options[2].value = "worst";
    downloadQualities.options[2].text = "worst"

  } else if (downloadMode == "Video And Audio") {
    document.getElementById("download-path").placeholder = "default: /playlists/%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s";
    if (input.length > 0) {
      let downloadPathRadio = document.getElementById("custom-download-output").checked;
      let youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
      if (downloadPathRadio == true) {
        document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s"
      } else if (youtubedlOutputRadio == true) {
        document.getElementById("output-path-indicator").innerHTML = input
      }
    } else {
      document.getElementById("output-path-indicator").innerHTML = "%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s"
    }

    fileExtensions.options[0].value = "any";
    fileExtensions.options[0].text = "any (recommended for now)";
    fileExtensions.options[1].value = "mp4";
    fileExtensions.options[1].text = "mp4";

    downloadQualities.options[0].value = "best";
    downloadQualities.options[0].text = "best";
    downloadQualities.options[1].value = "worst";
    downloadQualities.options[1].text = "worst"
  }
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
  let spinner = document.getElementById(id);
  spinner.classList.remove("d-none");
}

function stopSpinner(id) {
  let spinner = document.getElementById(id);
  spinner.classList.add("d-none")
}

function customYtdl(checkboxId) {
  document.getElementById("download-path").disabled = false
  if (checkboxId == "custom-download-output") {
    document.getElementById("download-path").placeholder = "default: /playlists/"
  } else if (checkboxId == "custom-ytdl-output") {
    document.getElementById("download-path").placeholder = "default: /playlists/%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s"
  }
}

function changeOutputPathIndicator(id) {
  document.getElementById("output-path-indicator").innerHTML = "";
  let downloadPathRadio = document.getElementById("custom-download-output").checked;
  let youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
  let input = document.getElementById(id).value;
  let downloadMode = document.getElementById("download-mode").value;
  if (downloadMode == "Audio Only") {
    if (downloadPathRadio == true) {
      document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s"
    } else if (youtubedlOutputRadio == true) {
      document.getElementById("output-path-indicator").innerHTML = input
    }
  } else if (downloadMode == "Video And Audio") {
    if (downloadPathRadio == true) {
      document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s"
    } else if (youtubedlOutputRadio == true) {
      document.getElementById("output-path-indicator").innerHTML = input
    }
  }
}
