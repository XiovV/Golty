let channels = [];
function checkAll() {
  fetch("http://localhost:8080/api/check-all")
  .then(res => res.json())
  .then(res => {
    if (res.Type == "Success") {
      displaySuccessMessage(res.Message)
    }
  })
}

function getChannels() {
  fetch("http://localhost:8080/api/get-channels")
  .then(res => res.json())
  .then(channels => {
    displayChannels(channels)
  })
}

function deleteChannel(id) {
  let channelURL = {
    "channelURL": id
  }

  const options = {
    method: 'POST',
    body: JSON.stringify(channelURL),
    headers: new Headers({
      'Content-Type': 'application/json'
    })
  }

  fetch("http://localhost:8080/api/delete-channel", options)
  .then(res => res.json())
  .then(res => {
    if (res.Type == "Success") {
      displaySuccessMessage(res.Message)

      removeFromList(channelURL.channelURL)
    }
  })
}

function removeFromList(channelURL) {
  channelURL = channelURL.replace('delChannel', '')

  let channels = document.getElementById("channels")
  let li = document.getElementById(channelURL + "listElem")
  console.log(channels)
  channels.removeChild(li)
}

function checkChannel(id) {
  let channelURL = {
    "channelURL": id
  }

  const options = {
    method: 'POST',
    body: JSON.stringify(channelURL),
    headers: new Headers({
      'Content-Type': 'application/json'
    })
  }

  fetch("http://localhost:8080/api/check-channel", options)
  .then(res => res.json())
  .then(res => {
    console.log(res)

    if (res.Type == "Success") {
      if (res.Key == "NO_NEW_VIDEOS") {
        displayWarningMessage(res.Message)
      } else if(res.Key == "NEW_VIDEO_DETECTED") {
        displaySuccessMessage(res.Message)
      }
    }
  })
}

function displaySuccessMessage(message) {
  let success = document.getElementById("success")
  success.classList.remove("d-none")
  success.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`
}

function displayWarningMessage(message) {
  let warning = document.getElementById("warning")
  warning.classList.remove("d-none")
  warning.innerHTML = `${message} <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>`
}

function displayChannels(channels) {
  let ul = document.getElementById("channels");
  ul.innerHTML = ""

  channels.forEach(channel => {
    let li = document.createElement("li");
    li.setAttribute("class", "list-group-item")
    li.setAttribute("id", channel.ChannelURL+"listElem")

    li.appendChild(document.createTextNode(channel.ChannelURL));
    ul.appendChild(li)
  })

  displayButtons()
}

function displayButtons() {
  channels = document.querySelectorAll(".list-group-item")
  channels.forEach(channel => {
    oldHTML = channel.innerHTML
    channel.innerHTML = oldHTML + `<button class="btn btn-danger float-right ml-2" id="${channel.innerHTML+"delChannel"}" onClick="deleteChannel(this.id)">&times</button><button class="btn btn-primary float-right" id="${channel.innerHTML}" onClick="checkChannel(this.id)">Check Now</button>`
  })
}