package api

import (
	"net/http"
	"strconv"

	"goFinal/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJsonError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if !isValidID(id) {
		writeJsonError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, "Task not found", http.StatusNotFound)
		return
	}

	writeJson(w, task)
}

func isValidID(id string) bool {
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt <= 0 {
		return false
	}
	return true
}