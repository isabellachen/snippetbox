package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"snippetbox.isachen.com/internal/models"
)

type snippetResponse struct {
	Result models.Snippet `json:"result"`
}

func validateParam(param string) (int, error) {
	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) HomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		app.notFound(w, http.StatusText(http.StatusNotFound))
	}

	baseTemplatePath := filepath.Join(app.basePath, "/ui/html/pages/base.tmpl.html")
	homeTemplatePath := filepath.Join(app.basePath, "/ui/html/pages/home.tmpl.html")
	navTemplatePath := filepath.Join(app.basePath, "/ui/html/partials/nav.tmpl.html")

	files := []string{
		baseTemplatePath,
		homeTemplatePath,
		navTemplatePath,
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) SnippetCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	item := &struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}

	if err := json.NewDecoder(req.Body).Decode(item); err != nil {
		message := fmt.Sprintf("Invalid JSON: %s", err)
		app.clientError(w, http.StatusMethodNotAllowed, message)
	}

	id, err := app.repo.Create(item.Title, item.Content)

	if err != nil {
		app.serverError(w, err)
	}

	message := fmt.Sprintf("Created a new snippet with id %d", id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}

func (app *application) SnippetView(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("id")
	id, err := validateParam(param)

	if err != nil {
		message := fmt.Sprintf("Item %s does not exist", param)
		app.notFound(w, message)
		return
	}

	snippet, err := app.repo.ById(id)
	if err != nil {
		app.notFound(w, err.Error())
	}

	response := &snippetResponse{
		Result: snippet,
	}

	b, err := json.Marshal(response)

	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
