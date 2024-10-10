package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Remove "/assets" prefix and clean the path
	fp, _ := strings.CutPrefix(filepath.Clean(r.URL.Path), "/assets")
	// Join path with the base directory
	// fmt.Println("directory : ", fp)
	fp = filepath.Join("../web/static", fp)

	// Check if request a file not a directory, and exist.
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, fp)
}
