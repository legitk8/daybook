package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"daybook/backend/services"
)

type PageHandler struct {
	service *services.PageService
}

func NewPageHandler(s *services.PageService) *PageHandler {
	return &PageHandler{
		service: s,
	}
}

func (h *PageHandler) CreatePageHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TEMP: Hardcoded user (until auth)
	userID := services.DevUserID

	page, err := h.service.CreatePage(
		r.Context(),
		userID,
		request.Title,
		request.Content,
	)
	if err != nil {
		log.Println("CreatePage failed:", err)
		http.Error(w, "Failed to create page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func (h *PageHandler) GetPageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseUUID(idStr, w)
	if !ok {
		return
	}

	page, err := h.service.GetPageByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func (h *PageHandler) UpdatePageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseUUID(idStr, w)
	if !ok {
		return
	}

	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	page, err := h.service.UpdatePage(
		r.Context(),
		id,
		request.Title,
		request.Content,
	)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func (h *PageHandler) ListPagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := services.DevUserID

	pages, err := h.service.ListPages(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to list pages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

func (h *PageHandler) DeletePageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseUUID(idStr, w)
	if !ok {
		return
	}

	if err := h.service.DeletePage(r.Context(), id); err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
