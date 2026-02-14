package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"daybook/backend/services"
)

func CreatePageHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	page := services.CreatePage(request.Title, request.Content)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func UpdatePageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedPage, updated := services.UpdatePage(id, request.Title, request.Content)

	if !updated {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPage)
}

func GetPageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	page, found := services.GetPageByID(id)
	if !found {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func ListPagesHandler(w http.ResponseWriter, r *http.Request) {
	pages := services.ListPages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

func DeletePageHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Page ID", http.StatusBadRequest)
		return
	}

	deleted := services.DeletePage(id)
	if !deleted {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
