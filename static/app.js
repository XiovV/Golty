let channels = [];
function checkAll() {
  fetch("http://localhost:8080/api/check-all")
  .then(res => res.json())
  .then(res => {
    console.log(res)
  })
}

function getChannels() {
  let ul = document.getElementById("channels");
  
  fetch("http://localhost:8080/api/get-channels")
  .then(res => res.json())
  .then(channels => {
    channels.forEach(channel => {
      let li = document.createElement("li");
      li.setAttribute("class", "list-group-item")

      li.appendChild(document.createTextNode(channel.ChannelURL));
      ul.appendChild(li)
    })

    displayButtons()
  })
}

function displayButtons() {
  channels = document.querySelectorAll(".list-group-item")
  channels.forEach(channel => {
    oldHTML = channel.innerHTML
    channel.innerHTML = oldHTML + `<button class="btn btn-danger float-right ml-2" id="${channel.innerHTML+"delChannel"}" onClick="deleteChannel(this.id)">&times</button><button class="btn btn-primary float-right" id="${channel.innerHTML}" onClick="checkChannel(this.id)">Check Now</button>`
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
    let success = document.getElementById("success");

    console.log(res)

    if (res.Type == "Success") {
      success.classList.remove("d-none")
      success.innerHTML = `${res.Message}`
    }
  })
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
    let success = document.getElementById("success");

    console.log(res)

    if (res.Type == "Success") {
      success.classList.remove("d-none")
      success.innerHTML = `${res.Message}`
    }
    
  })
  console.log("ID: ", id)
}