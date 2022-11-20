package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.isachen.com/internal/models"
)

type snippetResponse struct {
	Result *models.Snippet `json:"result"`
}

type latestSnippetsResponse struct {
	Result []*models.Snippet `json:"result"`
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

	snippets, err := app.repo.Latest(10)
	if err != nil {
		app.serverError(w, err)
	}

	data := app.newTemplateData(req)

	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
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

	id, err := app.repo.Create(item.Title, item.Content, 24)

	if err != nil {
		app.serverError(w, err)
	}

	message := fmt.Sprintf("Created a new snippet with id %d", id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}

func (app *application) SnippetView(w http.ResponseWriter, req *http.Request) {
	query := "id"
	param := req.URL.Query().Get(query)

	if param == "" {
		query = "limit"
		param = req.URL.Query().Get(query)
	}

	validatedParam, err := validateParam(param)

	if err != nil {
		message := fmt.Sprintf("Unable to validate param %s", param)
		app.notFound(w, message)
		return
	}

	if query == "id" {
		snippet, err := app.repo.ById(validatedParam)

		if err != nil {
			app.notFound(w, err.Error())
		}
		data := app.newTemplateData(req)
		data.Snippet = snippet

		app.render(w, http.StatusOK, "view.tmpl.html", data)
	}

	if query == "limit" {
		snippets, err := app.repo.Latest(validatedParam)
		if err != nil {
			app.serverError(w, err)
		}

		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(req)
		data.Snippets = snippets

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.serverError(w, err)
		}
	}
}
