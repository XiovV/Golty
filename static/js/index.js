let channels = [];

function addChannel() {
  startSpinner("add-channel-spinner");
  let downloadEntire = document.querySelector('#download-entire-channel').checked;
  let URL = document.getElementById("channel-url").value;
  let downloadMode = document.getElementById("download-mode").value;
  let fileExtension = document.getElementById("file-extension").value;
  let downloadQuality = document.getElementById("download-quality").value;
  let downloadPath = document.getElementById("output-path-indicator").innerText;

  console.log("DL PATH: ", downloadPath)

  let type = "Channel";

  let channelData = {
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
    body: JSON.stringify(channelData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/add", options)
    .then(res => res.json())
    .then(res => {
      handleResponse(res);
      stopSpinner("add-channel-spinner");
      getChannels()
    });
}

function checkAll() {
  startSpinner("check-all-spinner");
  let type = "channels";
  let channelData = {
    type,
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
        handleResponse(res);
        stopSpinner("check-all-spinner");
        getChannels()
      });
}

function getChannels() {
  let channelData = {
    Type: "channels"
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
      .then(channels => {
        displayChannels(channels);
        getVersion();
      });
}

function displayChannels(channels) {
  document.getElementById("accordion").innerHTML = "";
  console.log(channels);
  channels.forEach((channel, index) => {
    console.log(channel);
    document.getElementById("accordion").innerHTML += `<div class="mb-2 p-2 card">
      <h5 class="mb-0">
        <button class="btn btn-link dropdown-toggle" data-toggle="collapse" data-target="#collapse${index}" aria-expanded="true" aria-controls="collapse${index}" id=${channel.URL}listElem>
          ${channel.Name}
        </button><button class="btn btn-danger float-right ml-2" id="${channel.URL +
    "delTarget"}" onClick="deleteChannel(this.id)">&times</button><button class="btn btn-primary float-right" id="${channel.URL}" onClick="checkChannel(this.id)">Check<div id="${channel.URL}-spinner" class="spinner-border align-middle ml-2 d-none"></div></button>
      </h5>
  
      <div id="collapse${index}" class="collapse" aria-labelledby="heading${index}" data-parent="#accordion">
        <div class="panel-body ml-2">
          Latest Download: <a href=https://www.youtube.com/watch?v=${channel.LatestDownloaded} target="_blank">https://www.youtube.com/watch?v=${channel.LatestDownloaded}</a>
          <p>Download Mode: ${channel.DownloadMode}</p>
          <p>Last Checked: ${channel.LastChecked}</p>
          <p>Preferred Extension For Audio: ${channel.PreferredExtensionForAudio}
          <p>Preferred Extension For Video: ${channel.PreferredExtensionForVideo}
          <p class="m-0 p-0">Download Path: ${channel.DownloadPath}</p>
          <br>
          <button class="btn btn-link dropdown-toggle m-0 p-0" type="button" data-toggle="collapse" data-target="#history${index}" aria-expanded="false" aria-controls="history${index}">
            Download History
          </button>
          <div class="collapse m-0 p-0" id="history${index}">
            <div class="card card-body p-2" id="dlhistory${channel.Name}">
            </div>
          </div>
        </div>
      </div>
    </div>`;
    displayDownloadHistory(channel.Name, channel.DownloadHistory)
  })
}

function checkChannel(id) {
  startSpinner(id+"-spinner");
  let URL = id;
  let downloadMode = document.getElementById("download-mode").value;
  let fileExtension = document.getElementById("file-extension").value;
  let downloadQuality = document.getElementById("download-quality").value;
  let type = "Channel";

  let channelData = {
    URL,
    downloadMode,
    fileExtension,
    downloadQuality,
    type
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelData),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/check", options)
      .then(res => res.json())
      .then(res => {
        console.log(res);
        stopSpinner(id+"-spinner");
        if (res.Type === "Success") {
          if (res.Key === "NO_NEW_VIDEOS") {
            displayWarningMessage(res.Message);
            getChannels()
          } else if (res.Key === "NEW_VIDEO_DETECTED") {
            displaySuccessMessage(res.Message);
            getChannels()
          }
        } else if (res.Type === "Error") {
          if (res.Key === "ERROR_DOWNLOADING_VIDEO") {
            displayErrorMessage(res.Message);
          }
        }
      });
}

function deleteChannel(id) {
  let channelURL = {
    URL: id,
    Type: "Channel"
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelURL),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/delete", options)
    .then(res => res.json())
    .then(res => {
      handleResponse(res);
      getChannels()
    });
}

function displayWarningMessage(message) {
  let warning = document.getElementById("warning");
  warning.innerHTML = "";
  warning.classList.remove("d-none");
  warning.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function changeExtension() {
  console.log("change ext");

  let downloadMode = document.getElementById("download-mode").value;
  let fileExtensions = document.getElementById("file-extension");
  let input = document.getElementById("download-path").value;
  if (downloadMode === "Audio Only") {
    document.getElementById("download-path").placeholder = "default: /channels/%(uploader)s/audio/%(title)s.%(ext)s";
    if (input.length > 0) {
      let downloadPathRadio = document.getElementById("custom-download-output").checked;
      let youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
      if (downloadPathRadio === true) {
        document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/audio/%(title)s.%(ext)s"
      } else if (youtubedlOutputRadio === true) {
        document.getElementById("output-path-indicator").innerHTML = input
      }
    } else {
      document.getElementById("output-path-indicator").innerHTML = "/channels/%(uploader)s/audio/%(title)s.%(ext)s"
    }

    fileExtensions.options[0].value = "m4a";
    fileExtensions.options[0].text = "m4a (recommended for now)";
    fileExtensions.options[1].value = "mp3";
    fileExtensions.options[1].text = "mp3";
  } else if (downloadMode === "Video And Audio") {
    document.getElementById("download-path").placeholder = "default: /channels/%(uploader)s/video/%(title)s.%(ext)s";
    if (input.length > 0) {
      let downloadPathRadio = document.getElementById("custom-download-output").checked;
      let youtubedlOutputRadio = document.getElementById("custom-ytdl-output").checked;
      if (downloadPathRadio === true) {
        document.getElementById("output-path-indicator").innerHTML = input + "%(uploader)s/video/%(title)s.%(ext)s"
      } else if (youtubedlOutputRadio === true) {
        document.getElementById("output-path-indicator").innerHTML = input
      }
    } else {
      document.getElementById("output-path-indicator").innerHTML = "/channels/%(uploader)s/video/%(title)s.%(ext)s"
    }

    fileExtensions.options[0].value = "any";
    fileExtensions.options[0].text = "any (recommended for now)";
    fileExtensions.options[1].value = "mp4";
    fileExtensions.options[1].text = "mp4";
  }
}

function customYtdl(checkboxId) {
  document.getElementById("download-path").disabled = false;
  if (checkboxId === "custom-download-output") {
    document.getElementById("download-path").placeholder = "default: /channels/"
  } else if (checkboxId === "custom-ytdl-output") {
    document.getElementById("download-path").placeholder = "default: /channels/%(uploader)s/audio/%(title)s.%(ext)s"
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