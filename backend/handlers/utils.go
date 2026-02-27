package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func parseUUID(idStr string, w http.ResponseWriter) (uuid.UUID, bool) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return uuid.Nil, false
	}
	return id, true
}
