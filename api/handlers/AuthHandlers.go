package handlers

import (
	"GoServer/api/middleware"
	"GoServer/internal/db"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

// User struct to hold form data
type User struct {
	Username string
	Email    string
	Password string
}

// SignUpHandler renders the sign-up page and handles form submission
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmplPath := filepath.Join("web", "templates", "sign-up.html")
		renderTemplate(w, tmplPath, nil)
	} else if r.Method == http.MethodPost {
		user := User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// Validate form data
		errors := validateForm(user)
		if len(errors) > 0 {
			tmplPath := filepath.Join("web", "templates", "sign-up.html")
			data := struct {
				User   User
				Errors []string
			}{
				User:   user,
				Errors: errors,
			}
			renderTemplate(w, tmplPath, data)
			return
		}

		// Store user in the database
		err := storeUser(user)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			http.Error(w, "Error saving user to database", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
	}
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Regular expression for validating an email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validateForm validates the form data and returns a slice of error messages
func validateForm(user User) []string {
	var errors []string

	// Check if username is empty
	if user.Username == "" {
		errors = append(errors, "Username is required.")
	}

	// Check if email is empty
	if user.Email == "" {
		errors = append(errors, "Email is required.")
	} else {
		// Check if email format is valid
		if !isValidEmail(user.Email) {
			errors = append(errors, "Invalid email format.")
		}
	}

	// Check if password is empty
	if user.Password == "" {
		errors = append(errors, "Password is required.")
	} else {
		// Check if password length is at least 6 characters
		if len(user.Password) < 6 {
			errors = append(errors, "Password must be at least 6 characters long.")
		}
	}

	return errors
}

// LoginHandler renders the login page and handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmplPath := filepath.Join("web", "templates", "log-in.html")
		renderTemplate(w, tmplPath, nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Authenticate user
		if isValidUser(username, password) {
			// Create session
			sessionID := middleware.CreateSession(username)

			// Set a cookie with the session ID
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: sessionID,
				Path:  "/",
			})

			// Redirect to a dashboard or profile page
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		// Authentication failed, render login page with error
		tmplPath := filepath.Join("web", "templates", "log-in.html")
		data := struct {
			Error string
		}{
			Error: "Invalid username or password",
		}
		renderTemplate(w, tmplPath, data)
	}
}

// LogoutHandler handles log out request
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session (remove session_id cookie)
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Set MaxAge to -1 to delete the cookie
	})

	// Redirect to the home page or login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getUserByUsername retrieves user details from the database by username
func getUserByUsername(username string) (User, error) {
	var user User
	row := db.DB.QueryRow("SELECT username, email FROM users WHERE username = ?", username)
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		return user, errors.Wrap(err, "failed to retrieve user from database")
	}
	return user, nil
}

// storeUser stores a new user in the database
func storeUser(user User) error {
	_, err := db.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
	}
	return err
}

// isValidUser checks user credentials against the database
func isValidUser(username, password string) bool {
	var storedPassword string
	err := db.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err == sql.ErrNoRows || storedPassword != password {
		return false
	}
	return true
}
