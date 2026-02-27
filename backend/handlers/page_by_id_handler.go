package handlers

import (
	"net/http"
	"strings"
)

func (h *PageHandler) PageByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/pages/")

	switch r.Method {
	case http.MethodGet:
		h.GetPageHandler(w, r, id)
	case http.MethodPut:
		h.UpdatePageHandler(w, r, id)
	case http.MethodDelete:
		h.DeletePageHandler(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
