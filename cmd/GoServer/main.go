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

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define your handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/sign-up", handlers.SignUpHandler)
	http.HandleFunc("/log-in", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/log-out", handlers.LogoutHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	// http.HandleFunc("/music-playlist", handlers.MusicHandler)

	http.HandleFunc("/store", handlers.RenderStore)
	http.HandleFunc("/add-to-cart", handlers.AddToCartHandler)
	http.HandleFunc("/clear-cart", handlers.ClearCartHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}

}
