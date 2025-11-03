package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// Directory to store static files
	staticDir := "./static"

	// Create a file server handler - picks file from the local
	fileServer := http.FileServer(http.Dir(staticDir))

	// Tell the server to serve everything under /static/
	// Handle - when you already have a handler
	// It can take custom handlers as well
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Optional: serve index.html when you visit "/"
	// Handle - custom handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	// Define port
	port := "8080"
	fmt.Printf("Server started at http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
