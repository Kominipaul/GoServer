// cmd/GoServer/main.go

package main

import (
	"log"
	"net/http"
	"os"

	"GoServer/api/handlers"
	"GoServer/internal/db"
)

func main() {
	// Initialize the database connection
	db.Init()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
