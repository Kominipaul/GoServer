package main

import (
	"log"
	"net/http"

	"GoServer/api/handlers" // Absolute import path
	// Absolute import path
)

func main() {
	// Register handlers and middleware
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/sign-up", handlers.SignUpHandler)
	http.HandleFunc("/log-in", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/log-out", handlers.LogoutHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Println("Starting server on :8081")
	log.Println("URL http://localhost:8081")

	// Start the server
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
