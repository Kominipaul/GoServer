// GoServer/api/handlers/handlers.go

package handlers

import (
	"GoServer/api/middleware"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/pkg/errors"
)

// User struct to hold form data
type User struct {
	Username string
	Email    string
	Password string
}

var db *sql.DB

func init() {
	// Connect to the MySQL database
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/myapp")
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to database"))
	}
}

// HomeHandler renders the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	renderTemplate(w, tmplPath, nil)
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

		// Insert user into database
		err := insertUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the home page or login page
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
	}
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

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Regular expression for validating an email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// insertUser inserts a new user into the database
func insertUser(user User) error {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		return errors.Wrap(err, "failed to insert user into database")
	}
	return nil
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
			// Generate a unique session ID (for demonstration, using username)
			sessionID := username

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

// DashboardHandler handles dashboard page (protected route)
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
		return
	}

	// Get username from session (for demonstration, using username as session ID)
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Query user details from database
	user, err := getUserByUsername(sessionID.Value)
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
	}{
		Username: user.Username,
	}

	// Execute template
	tmpl.Execute(w, data)
}

// getUserByUsername retrieves user details from the database by username
func getUserByUsername(username string) (User, error) {
	var user User
	row := db.QueryRow("SELECT username, email FROM users WHERE username = ?", username)
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		return user, errors.Wrap(err, "failed to retrieve user from database")
	}
	return user, nil
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

// isValidUser checks if the provided username and password are valid
func isValidUser(username, password string) bool {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		return false
	}
	return storedPassword == password
}
