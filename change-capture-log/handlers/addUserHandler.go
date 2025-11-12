package handlers

import (
	main "changecapturelog"
	"net/http"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
	}

	user := new(main.Users)
}
