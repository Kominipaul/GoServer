package main

import (
    "log"
    "net/http"
    "GoServer/api/handlers"
)

func main() {
    // This func is responsible for routing the requests to the appropriate handlers
    // The first argument is the route, the second is the handler function
    http.HandleFunc("/", handlers.HomeHandler)

    // This is a simple way to serve static files
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

    log.Println("Starting server on :8080")
    log.Println("URL http://localhost:8080")

    // Start the server
    err := http.ListenAndServe(":8080", nil)

    // If there was an error starting the server, log the error and exit
    // This is a fatal error, so the program will exit
    // This can hepend if the port is already in use or the user does not have permission to use it
    if err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}
