package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func validateParam(param string) error {
	_, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	return nil
}

func (application *application) HomeHandler(w http.ResponseWriter, res *http.Request) {
	if res.URL.Path != "/" {
		application.notFound(w, http.StatusText(http.StatusNotFound))
	}

	files := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		application.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		application.serverError(w, err)
	}
	w.Write([]byte("hello from snippetbox"))
}

func (application *application) SnippetCreate(w http.ResponseWriter, res *http.Request) {
	if res.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		application.clientError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func (application *application) SnippetView(w http.ResponseWriter, res *http.Request) {
	param := res.URL.Query().Get("id")
	if err := validateParam(param); err != nil {
		message := fmt.Sprintf("Item %s does not exist", param)
		application.notFound(w, message)

		return
	}
	message := fmt.Sprintf("Snippet %s", param)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))
}

func newMux(cfg *config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.static))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", cfg.Applicaion.HomeHandler)
	mux.HandleFunc("/snippet/create", cfg.Applicaion.SnippetCreate)
	mux.HandleFunc("/snippet/view", cfg.Applicaion.SnippetView)
	return mux
}
