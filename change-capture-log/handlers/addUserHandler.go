package handlers

import "net/http"

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
	}
}
