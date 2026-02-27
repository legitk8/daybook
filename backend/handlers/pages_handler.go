package handlers

import "net/http"

func (h *PageHandler) PagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListPagesHandler(w, r)
	case http.MethodPost:
		h.CreatePageHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
