//go:build !prod

package api

import (
	"fmt"
	"net/http"
	"os"
)

func getFrontendAssets() http.FileSystem {
	fmt.Println("using dev mode")
	return http.FS(os.DirFS("dist"))
}
