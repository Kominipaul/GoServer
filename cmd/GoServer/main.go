// cmd/GoServer/main.go

package main

import (
	"GoServer/api/handlers"
	"GoServer/internal/db"
	"log"
	"net/http"
	"os"
)

func main() {

	// Initialize the database connection
	db.Init()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/sign-up", handlers.SignUpHandler)
	http.HandleFunc("/log-in", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/log-out", handlers.LogoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}

}
