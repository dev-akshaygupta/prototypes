package handlers

import "net/http"

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
	}
}
