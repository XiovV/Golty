let channels = [];

function addChannel() {
  startSpinner("add-channel-spinner")

  let channelURL = document.getElementById("channel-url").value
  let downloadMode = document.getElementById("download-mode").value
  let fileExtension = document.getElementById("file-extension").value
  let downloadQuality = document.getElementById("download-quality").value

  let channelData = {
    channelURL,
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

  fetch("/api/add-channel", options)
    .then(res => res.json())
    .then(res => {
      if (res.Type == "Success") {
        displaySuccessMessage(res.Message);
        stopSpinner("add-channel-spinner")
        getChannels()
      }
    });
}

function checkAll() {
  fetch("/api/check-all")
    .then(res => res.json())
    .then(res => {
      if (res.Type == "Success") {
        displaySuccessMessage(res.Message);
      }
    });
}

function getChannels() {
  fetch("/api/get-channels")
    .then(res => res.json())
    .then(channels => {
      displayChannels(channels);
    });
}

function deleteChannel(id) {
  let channelURL = {
    channelURL: id
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelURL),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/delete-channel", options)
    .then(res => res.json())
    .then(res => {
      if (res.Type == "Success") {
        displaySuccessMessage(res.Message);

        removeFromList(channelURL.channelURL);
      }
    });
}

function removeFromList(channelURL) {
  channelURL = channelURL.replace("delChannel", "");

  let channels = document.getElementById("channels");
  let li = document.getElementById(channelURL + "listElem");
  console.log(channels);
  channels.removeChild(li);
}

function checkChannel(id) {
  let channelURL = {
    channelURL: id
  };

  const options = {
    method: "POST",
    body: JSON.stringify(channelURL),
    headers: new Headers({
      "Content-Type": "application/json"
    })
  };

  fetch("/api/check-channel", options)
    .then(res => res.json())
    .then(res => {
      console.log(res);

      if (res.Type == "Success") {
        if (res.Key == "NO_NEW_VIDEOS") {
          displayWarningMessage(res.Message);
        } else if (res.Key == "NEW_VIDEO_DETECTED") {
          displaySuccessMessage(res.Message);
        }
      } else if (res.Type == "Error") {
        if (res.Key == "ERROR_DOWNLOADING_VIDEO") {
          displayErrorMessage(res.Message);
        }
      }
    });
}
function displayErrorMessage(message) {
  let error = document.getElementById("error");
  error.classList.remove("d-none");
  error.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displaySuccessMessage(message) {
  let success = document.getElementById("success");
  success.classList.remove("d-none");
  success.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displayWarningMessage(message) {
  let warning = document.getElementById("warning");
  warning.classList.remove("d-none");
  warning.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`;
}

function displayChannels(channels) {
  console.log(channels)
  channels.forEach((channel, index) => {
    console.log(channel)
    document.getElementById("accordion").innerHTML += `<div class="card">
      <h5 class="mb-0">
        <button class="btn btn-link" data-toggle="collapse" data-target="#collapse${index}" aria-expanded="true" aria-controls="collapse${index}" id=${channel.ChannelURL}listElem>
          ${channel.ChannelURL}
        </button><button class="btn btn-danger float-right ml-2" id="${channel.innerHTML +
          "delChannel"}" onClick="deleteChannel(this.id)">&times</button><button class="btn btn-primary float-right" id="${
          channel.innerHTML
        }" onClick="checkChannel(this.id)">Check Now</button>
      </h5>
  
    <div id="collapse${index}" class="collapse" aria-labelledby="heading${index}" data-parent="#accordion">
      <div class="card-body">
        Basic channel info
      </div>
    </div>
  </div>`
  })

  displayButtons();
}

function displayButtons() {
  channels = document.querySelectorAll(".list-group-item");
  channels.forEach(channel => {
    oldHTML = channel.innerHTML;
    channel.innerHTML =
      oldHTML +
      `<button class="btn btn-danger float-right ml-2" id="${channel.innerHTML +
        "delChannel"}" onClick="deleteChannel(this.id)">&times</button><button class="btn btn-primary float-right" id="${
        channel.innerHTML
      }" onClick="checkChannel(this.id)">Check Now</button>`;
  });
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
    fileExtensions.options[0].value = ".m4a"
    fileExtensions.options[0].text = ".m4a"
    fileExtensions.options[1].value = ".mp3"
    fileExtensions.options[1].text = ".mp3"
    downloadQualities.options[0].value = "best"
    downloadQualities.options[0].text = "best"
    downloadQualities.options[1].value = "medium"
    downloadQualities.options[1].text = "medium"
    downloadQualities.options[2].value = "worst"
    downloadQualities.options[2].text = "worst"
  } else if (downloadMode == "Video And Audio") {
    fileExtensions.options[0].value = ".mkv"
    fileExtensions.options[0].text = ".mkv"
    fileExtensions.options[1].value = ".mp4"
    fileExtensions.options[1].text = ".mp4"
    downloadQualities.options[0].value = "best"
    downloadQualities.options[0].text = "best"
    downloadQualities.options[1] = null
    downloadQualities.options[2].value = "worst"
    downloadQualities.options[2].text = "worst"
  }
}