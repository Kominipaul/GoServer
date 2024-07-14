package handlers

import (
	"GoServer/api/middleware"
	"html/template"
	"net/http"
	"path/filepath"
)

// DashboardHandler handles dashboard page (protected route)
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
		return
	}

	// Get username from session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Use SessionsMutex and Sessions from middleware package
	middleware.SessionsMutex.Lock()
	defer middleware.SessionsMutex.Unlock()

	session, exists := middleware.Sessions[sessionID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Query user details from database based on session.Username
	user, err := getUserByUsername(session.Username)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Render dashboard template with user details
	tmplPath := filepath.Join("web", "templates", "dashboard.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass user data to template
	data := struct {
		Username string
		Email    string
	}{
		Username: user.Username,
		Email:    user.Email,
	}

	// Execute template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
