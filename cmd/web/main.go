package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Import with blank identifier so the library can
	// register itself with the database/sql package
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.isachen.com/internal/models"
	"snippetbox.isachen.com/internal/repository"
)

type config struct {
	addr   string
	static string
	dsn    string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	basePath string
	repo     models.Repository
}

func main() {
	cwd, _ := os.Getwd()
	repo := repository.NewInMemoryRepo()

	app := &application{
		repo:     repo,
		basePath: cwd,
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	cfg := &config{}

	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.addr, "addr", "8080", "HTTP network address")
	flag.StringVar(&cfg.static, "static", "./ui/static/", "Path for static files")

	// Call flag.Parse() only after all the flags have been declared
	flag.Parse()
	db, err := openDB(cfg.dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.infoLog.Printf("Listening on port %v...", cfg.addr)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", "localhost", cfg.addr),
		Handler:      app.routes(cfg),
		ErrorLog:     app.errorLog,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		app.errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// Init a pool of several connections without connecting
	// Uses a imported library driver specific helper
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
