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

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// New dynamic middleware for dynamic routes, e.g. session manager
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// We don't need dynamic middleware for static files
	staticDir := os.Getenv("STATIC_DIR")
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
