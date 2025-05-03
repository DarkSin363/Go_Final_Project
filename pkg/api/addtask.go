package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"goFinal/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

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
		writeJsonError(w, "Date error", http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]interface{}{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = now.Format(F)
		return nil
	}

	t, err := time.Parse(F, task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(F)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
	}
}

func writeJsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Failed to encode error: %v", err)
	}
}
