//go:build prod

package api

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

var embedFrontend embed.FS

func getFrontendAssets() http.FileSystem {
	fmt.Println("using prod mode")
	fileSystem, err := fs.Sub(embedFrontend, "dist")
	if err != nil {
		log.Fatalln(err)
	}

	return http.FS(fileSystem)
}
