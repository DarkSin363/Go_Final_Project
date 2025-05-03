package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"goFinal/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")

	allTasks, err := db.GetAllTasks()
	if err != nil {
		writeJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	var result []*db.Task

	if search == "" {
		result = allTasks
	} else {
		result = filterTasks(allTasks, search)
	}

	log.Printf("Returning %d tasks for search '%s'", len(result), search)
	if result == nil {
		result = make([]*db.Task, 0)
	}

	writeJson(w, TasksResp{Tasks: result})
}

func filterTasks(tasks []*db.Task, search string) []*db.Task {
	var result []*db.Task
	searchLower := strings.ToLower(search)

	if date, err := time.Parse("02.01.2006", search); err == nil {
		searchDate := date.Format("20060102")
		for _, task := range tasks {
			if task.Date == searchDate {
				result = append(result, task)
			}
		}
		return result
	}

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), searchLower) ||
			strings.Contains(strings.ToLower(task.Comment), searchLower) {
			result = append(result, task)
		}
	}

	return result
}
