function getVersion() {
    let version = document.getElementById("version-number")

    fetch("/api/version")
        .then(res => res.json())
        .then(res => {
            version.innerText = res.VersionNumber
        })
}