package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// HomeHandler renders the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	renderTemplate(w, tmplPath, nil)
}

// renderTemplate parses and executes the specified template file
func renderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

