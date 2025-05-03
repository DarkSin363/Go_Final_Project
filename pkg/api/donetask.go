package api

import (
	"net/http"
	"time"

	"goFinal/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJsonError(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJsonError(w, "Date error", http.StatusBadRequest)
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJsonError(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, map[string]interface{}{})
}