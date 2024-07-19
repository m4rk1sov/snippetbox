package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// .env variables initialization
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	netAddress := os.Getenv("NET_ADDR")
	staticDir := os.Getenv("STATIC_DIR")

	// use .env variables | .String to convert user input of flag
	addr := flag.String("addr", netAddress, "http network address")
	if addr == nil {
		log.Fatalf("missing env var: %s", netAddress)
	}

	// This command is responsible for taking command line value and assigning it
	// to a values for flag, call before usage of variables of flag.
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// path to static assets directory in .env file
	fileServer := http.FileServer(http.Dir(staticDir))
	if fileServer == nil {
		log.Fatal("Error creating file server")
	}
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Returning the address of addr variable
	log.Printf("Starting the server on %s", *addr)
	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
