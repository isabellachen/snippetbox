package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var cfg config

type config struct {
	addr   string
	static string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	app := &application{}

	flag.StringVar(&cfg.addr, "addr", "8080", "HTTP network address")
	flag.StringVar(&cfg.static, "static", "./ui/static/", "Path for static files")

	// Call flag.Parse() only after all the flags have been declared
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.infoLog = infoLog
	app.errorLog = errorLog

	infoLog.Printf("Listening on port %v...", cfg.addr)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", "localhost", cfg.addr),
		Handler:      newMux(app),
		ErrorLog:     errorLog,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
