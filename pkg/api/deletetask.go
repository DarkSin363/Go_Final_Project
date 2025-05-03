package api

import (
	"net/http"

	"goFinal/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJsonError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if !isValidID(id) {
		writeJsonError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]interface{}{})
}
