package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
	http.Handler - any type that has method - ServeHTTP(w http.ResponseWriter, r *http.Request)
	ServeHTTP(w http.ResponseWriter, r *http.Request) - interface
	it can be added to structs or functions
	struct is like a custom object which translates to object of handler (in leyman terms)
*/
// custom handler - HelloHandler
type GreetingHandler struct {
	Message string
}

func (g GreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, g.Message)
	fmt.Fprintln(w, "This is written inside handler")
}

func main() {

	// Call to custome Handler - HelloHandler
	http.Handle("/customhandler", GreetingHandler{Message: "This message is passed into the handler while calling it!!"})

	// Custom message - text
	// w http.ResponseWriter - consider like a notepad -> that browser will read from
	// "/about" - URL pattern it will respond to
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello this is a custom about page from HandleFunc")
	})

	// Custom message - HTML tags/page
	// w http.ResponseWriter - consider like a notepad -> that browser will read from
	http.HandleFunc("/htmlmsg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/HTML")
		fmt.Fprintln(w, "<h3>HTML from HandleFunc<h3>")
	})

	// Custom message - input from URL (query parameter)
	http.HandleFunc("/inputlink", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if name == "" {
			name = "Guest"
		}

		fmt.Fprintf(w, "Pleasure to meet you %s\n", name)
		fmt.Fprintf(w, "We have successfully received request for -  %s\n", name)
	})

	// Custom message - HTML form (POST request)
	http.HandleFunc("/thisform", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			message := r.FormValue("message")
			fmt.Fprintf(w, "Decoded message %s\n", message)
		} else {
			fmt.Fprintln(w, `<html><body><form method="POST" action="/thisform">
							<input name="message" placeholder="type something">
							<input type="submit">
							</form></html></body>`)
		}
	})

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
