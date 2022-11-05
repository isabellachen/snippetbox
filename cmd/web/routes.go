package main

import "net/http"

func (app *application) routes(cfg *config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.static))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	return mux
}
