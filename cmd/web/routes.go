package main

import (
	"github.com/bmizerany/pat"
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

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
