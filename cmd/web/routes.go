package main

import (
	"github.com/justinas/alice"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {

	// Was called in the main function
	//err := godotenv.Load(".env")
	//if err != nil {
	//	app.errorLog.Fatalf("Error loading .env file: %s", err)
	//}

	staticDir := os.Getenv("STATIC_DIR")

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
