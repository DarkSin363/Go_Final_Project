package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goFinal/pkg/db"
)

func putTaskHandler(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		writeJsonError(w, "Request body too large or corrupted", http.StatusBadRequest)
		return
	}

	if !json.Valid(buf.Bytes()) {
		writeJsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var task *db.Task
	log.Printf("Task : %v", task)
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJsonError(w, fmt.Sprintf("JSON parse error: %v", err), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJsonError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(task); err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.UpdateTask(task)
	if err != nil {
		writeJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]interface{}{})
}
