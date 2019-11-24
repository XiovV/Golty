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
        console.log("CREATE SUCCESS ALERT");
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

function displayDownloadHistory(channelName, downloadHistory) {
    let historyBox = document.getElementById("dlhistory"+channelName)
    console.log(historyBox)
    console.log("DISPLAY HISTORY")
    downloadHistory.forEach(video => {
        historyBox.innerHTML += `<br> <a href=https://www.youtube.com/watch?v=${video} target="_blank">https://www.youtube.com/watch?v=${video}</a>`
    })
}