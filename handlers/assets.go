package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

// parse all tamplates at once in the beggining of the program
func init() {
	templates = template.Must(templates.ParseGlob("./web/templates/*.html"))
	templates = template.Must(templates.ParseGlob("./web/templates/components/*.html"))
}

// AssetsHandler serves static files
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	fp, _ := strings.CutPrefix(r.URL.Path, "/assets")
	fp = filepath.Join("web/static", fp)

	_, err := os.Stat(fp)
	if err != nil || strings.HasSuffix(r.URL.Path, "/") {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, fp)
}
