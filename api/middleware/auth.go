// middleware/auth.go

package middleware

import (
	"net/http"
)

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	// Simplified check for demonstration (replace with proper session management)
	return cookie.Value != ""
}
