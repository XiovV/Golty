function changeExtension() {
    console.log("change ext")

    let downloadMode = document.getElementById("download-mode").value
    let fileExtensions = document.getElementById("file-extension")
    let downloadQualities = document.getElementById("download-quality")
    if (downloadMode == "Audio Only") {
        document.getElementById("download-path").placeholder = "default: videos/%(uploader)s/audio/%(title)s.%(ext)s"

        fileExtensions.options[0].value = "m4a"
        fileExtensions.options[0].text = "m4a"
        fileExtensions.options[1].value = "mp3"
        fileExtensions.options[1].text = "mp3"
        downloadQualities.options[0].value = "best"
        downloadQualities.options[0].text = "best"
        downloadQualities.options[1].value = "medium"
        downloadQualities.options[1].text = "medium"
        downloadQualities.options[2].value = "worst"
        downloadQualities.options[2].text = "worst"

    } else if (downloadMode == "Video And Audio") {
        document.getElementById("download-path").placeholder = "default: videos/%(uploader)s/video/%(title)s.%(ext)s"

        fileExtensions.options[0].value = "any"
        fileExtensions.options[0].text = "any (recommended for now)"
        fileExtensions.options[1].value = "mp4"
        fileExtensions.options[1].text = "mp4"
        // fileExtensions.options[2].value = ".mkv"
        // fileExtensions.options[2].text = ".mkv"

        downloadQualities.options[0].value = "best"
        downloadQualities.options[0].text = "best"
        // downloadQualities.options[1] = null
        downloadQualities.options[1].value = "worst"
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
    spinner = document.getElementById(id);
    spinner.classList.remove("d-none");
}

function stopSpinner(id) {
    spinner = document.getElementById(id);
    spinner.classList.add("d-none")
}
