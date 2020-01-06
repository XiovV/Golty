function handleResponse(res) {
    if (res.Type === "Success") {
        displaySuccessMessage(res.Message)
    } else if (res.Type === "Error") {
        displayErrorMessage(res.Message)
    } else if (res.Type === "Warning") {
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
    downloadHistory.forEach(video => {
        historyBox.innerHTML += `<br> <a href=https://www.youtube.com/watch?v=${video} target="_blank">https://www.youtube.com/watch?v=${video}</a>`
    })
}

function getCheckingInterval() {
    let type = document.getElementById("list-type").value;

    let getCheckingInterval = {
        type
    };

    const options = {
        method: "POST",
        body: JSON.stringify(getCheckingInterval),
        headers: new Headers({
            "Content-Type": "application/json"
        })
    };

    fetch("/api/get-checking-interval", options)
        .then(res => res.json())
        .then(interval => {
            document.getElementById("checking-interval").value = interval.CheckingInterval
            document.getElementById("time").value = interval.Time
        });
}

function updateCheckingInterval() {
    startSpinner("update-checking-interval-spinner");
    let checkingInterval = document.getElementById("checking-interval").value;
    let time = document.getElementById("time").value;
    let type = document.getElementById("list-type").value;

    let interval = {
      checkingInterval: checkingInterval.toString(),
      time: time,
      type
    };

    const options = {
      method: "POST",
      body: JSON.stringify(interval),
      headers: new Headers({
        "Content-Type": "application/json"
      })
    };

    fetch("/api/update-checking-interval", options)
        .then(res => res.json())
        .then(res => {
          handleResponse(res);
          stopSpinner("update-checking-interval-spinner")
        });
}
