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

function setCheckingInterval() {
    let checkingIntervalInput = document.getElementById("checking-interval").value;
    let time = document.getElementById("time").value;

    if (time === "minutes") {
        return checkingIntervalInput
    } else if (time === "hours") {
        return  checkingIntervalInput * 60
    } else if (time === "days") {
        return checkingIntervalInput * 1440
    }
}
