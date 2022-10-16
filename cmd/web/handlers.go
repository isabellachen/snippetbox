package main

import (
	"fmt"
	"html/template"
	"log"
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
	application.infoLog.Printf("Home handler called")
	if res.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	files := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// application.errorLog.Fatal(err.Error())
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {

		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	w.Write([]byte("hello from snippetbox"))
}

func (application *application) SnippetCreate(w http.ResponseWriter, res *http.Request) {
	if res.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func (application *application) SnippetView(w http.ResponseWriter, res *http.Request) {
	param := res.URL.Query().Get("id")
	if err := validateParam(param); err != nil {
		message := fmt.Sprintf("Item %s does not exist", param)
		http.Error(w, message, http.StatusNotFound)
		return
	}
	message := fmt.Sprintf("Snippet %s", param)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))
}

func newMux(app *application) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.static))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	return mux
}
