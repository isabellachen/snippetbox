package main

import "net/http"

func (app *application) routes(cfg *config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.static))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", cfg.Applicaion.HomeHandler)
	mux.HandleFunc("/snippet/create", cfg.Applicaion.SnippetCreate)
	mux.HandleFunc("/snippet/view", cfg.Applicaion.SnippetView)
	return mux
}
