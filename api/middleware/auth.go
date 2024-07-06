package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Session struct to hold session data
type Session struct {
	Username string
	Expires  time.Time
}

var (
	Sessions      = map[string]Session{} // map to store sessions
	SessionsMutex = &sync.Mutex{}        // mutex for concurrent access to sessions
)

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	sessionID := cookie.Value

	SessionsMutex.Lock()
	defer SessionsMutex.Unlock()

	session, exists := Sessions[sessionID]
	if !exists {
		return false
	}

	// Check if the session has expired
	if session.Expires.Before(time.Now()) {
		delete(Sessions, sessionID)
		return false
	}

	return true
}

// CreateSession creates a new session for the given username and returns the session ID
func CreateSession(username string) string {
	sessionID := generateSessionID()
	expires := time.Now().Add(24 * time.Hour) // Session valid for 24 hours

	SessionsMutex.Lock()
	defer SessionsMutex.Unlock()

	Sessions[sessionID] = Session{Username: username, Expires: expires}

	return sessionID
}

// generateSessionID generates a unique session ID (simplified for demonstration)
func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// LogoutUser invalidates the session for the given session ID
func LogoutUser(sessionID string) {
	SessionsMutex.Lock()
	defer SessionsMutex.Unlock()

	delete(Sessions, sessionID)
}
