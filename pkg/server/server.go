package server

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"strings"

	"goFinal/pkg/api"
	"goFinal/pkg/db"
	"github.com/joho/godotenv"
)

func Run() error {

	fmt.Println("Start server..")
    
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pathDB := os.Getenv("TODO_DBFILE")
	if pathDB == "" {
		pathDB = "scheduler.db"
	}
	if err := db.Init(pathDB); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	defer db.Close()

	port := ":7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = ":" + envPort
	}

	webDir := "web"
	apiRouter := api.Init()
	fs := http.FileServer(http.Dir(webDir))

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			apiRouter.ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	return http.ListenAndServe(port, mainHandler)
}
