package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

// HomeHandler is a simple handler function that writes a ResponseWriter
// This function is responsible for rendering the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // This is the path to the template filepath
    // This is a relative path, so it will be relative to the location of the executable
    tmplPath := filepath.Join("web", "templates", "index.html")
    // Parse the template file
    tmpl, err := template.ParseFiles(tmplPath)
    // If there was an error parsing the template, write an error to the ResponseWriter
    if err != nil {
        // This is a 500 error, so we write a 500 status code
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Execute the template, passing in nil as the data
    tmpl.Execute(w, nil)
}
