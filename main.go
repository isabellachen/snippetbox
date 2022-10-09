package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := 8080
	fmt.Printf("Listening on port %d...", port)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "localhost", port),
		Handler:      newMux(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
