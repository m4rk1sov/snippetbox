package main

import (
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"m4rk1sov/snippetbox/pkg/models/mysql"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

// Generally avoid using Panic() and Fatal() outside of main function, better return errors from functions
func main() {

	// Writing info logs to a file
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Initializing two new loggers, concurrently safe but needs Write() to be safe for concurrency as well
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// .env variables initialization
	err = godotenv.Load(".env")
	if err != nil {
		errorLog.Fatalf("Error loading .env file: %s", err)
	}
	netAddress := os.Getenv("NET_ADDR")

	// use .env variables | .String to convert user input of flag
	addr := flag.String("addr", netAddress, "http network address")
	if addr == nil {
		errorLog.Fatalf("missing environmental variable: %s", netAddress)
	}

	mysqlAddress := os.Getenv("SNIPPETBOX_DB_DSN")
	dsn := flag.String("dsn", mysqlAddress, "MySQL data source name")

	// Secret key, 32 bytes, for authentications of user sessions/cookies
	secretKey := os.Getenv("SECRET_KEY")
	secret := flag.String("secret", secretKey, "Secret key")

	// This command is responsible for taking command line value and assigning it
	// to a values for flag, call before usage of variables of flag.
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateDir := os.Getenv("TEMPLATE_DIR")
	templateCache, err := newTemplateCache(templateDir)
	if err != nil {
		errorLog.Fatal(err)
	}

	// New session manager, passing key as parameter, we can add different configuration attributes/settings
	// As an example, here we set the 12-hour expiry time
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// This http.Server struct uses previous variables and custom errorLog logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Returning the address of addr variable
	infoLog.Printf("Starting the server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
