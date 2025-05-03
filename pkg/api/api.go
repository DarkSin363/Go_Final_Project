package api

import (
	"github.com/go-chi/chi/v5"
)

func Init() *chi.Mux {
    r := chi.NewRouter()

    r.Post("/api/signin", signinHandler)
	r.Get("/api/nextdate", nextDayHandler)

    r.Group(func(protected chi.Router) {
        protected.Use(auth)

        protected.Get("/api/task", getTaskHandler)
        protected.Post("/api/task", addTaskHandler)
        protected.Put("/api/task", putTaskHandler)
        protected.Delete("/api/task", deleteTaskHandler)
        protected.Get("/api/tasks", tasksHandler)
        protected.Post("/api/task/done", doneTaskHandler)
    })

    return r
}