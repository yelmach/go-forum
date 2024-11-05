package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	fp, _ := strings.CutPrefix(filepath.Clean(r.URL.Path), "/assets")
	fp = filepath.Join("web/static", fp)

	_, err := os.Stat(fp)
	if err != nil || strings.HasSuffix(r.URL.Path, "/") {
		// ErrorHandler(w, r, http.StatusNotFound)
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, fp)
}
