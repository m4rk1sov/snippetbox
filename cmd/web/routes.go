package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {

	err := godotenv.Load(".env")
	if err != nil {
		app.errorLog.Fatalf("Error loading .env file: %s", err)
	}
	staticDir := os.Getenv("STATIC_DIR")

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
