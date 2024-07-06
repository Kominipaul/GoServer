// GoServer/api/middleware/auth.go

package middleware

import (
	"net/http"
	"time"
	"sync"
    "fmt"
)

// Session struct to hold session data
type Session struct {
	Username string
	Expires  time.Time
}

var (
	sessions      = map[string]Session{}
	sessionsMutex = &sync.Mutex{}
)

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	sessionID := cookie.Value

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, exists := sessions[sessionID]
	if !exists {
		return false
	}

	// Check if the session has expired
	if session.Expires.Before(time.Now()) {
		delete(sessions, sessionID)
		return false
	}

	return true
}

// CreateSession creates a new session for the given username and returns the session ID
func CreateSession(username string) string {
	sessionID := generateSessionID()
	expires := time.Now().Add(24 * time.Hour) // Session valid for 24 hours

	sessionsMutex.Lock()
	sessions[sessionID] = Session{Username: username, Expires: expires}
	sessionsMutex.Unlock()

	return sessionID
}

// generateSessionID generates a unique session ID (simplified for demonstration)
func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// LogoutUser invalidates the session for the given session ID
func LogoutUser(sessionID string) {
	sessionsMutex.Lock()
	delete(sessions, sessionID)
	sessionsMutex.Unlock()
}
