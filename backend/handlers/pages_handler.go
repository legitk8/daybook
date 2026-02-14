package handlers

import "net/http"

func PagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ListPagesHandler(w, r)

	case http.MethodPost:
		CreatePageHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
