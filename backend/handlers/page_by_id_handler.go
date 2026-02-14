package handlers

import (
	"net/http"
	"strings"
)

func PageByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/pages/")

	switch r.Method {
	case http.MethodGet:
		GetPageHandler(w, r, id)
	case http.MethodDelete:
		DeletePageHandler(w, r, id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
