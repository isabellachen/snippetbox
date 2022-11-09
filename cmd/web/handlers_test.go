package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"snippetbox.isachen.com/internal/repository"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()

	cwd, _ := os.Getwd()
	basePath := filepath.Join(cwd, "../..")

	repo := repository.NewInMemoryRepo()

	app := &application{
		repo:     repo,
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		basePath: basePath,
	}

	cfg := &config{}

	testServer := httptest.NewServer(app.routes(cfg))

	return testServer.URL, func() {
		testServer.Close()
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		path               string
		expectedStatusCode int
		expectedContent    string
	}{
		{path: "/",
			expectedStatusCode: http.StatusOK,
			expectedContent:    "Home - Snippetbox",
		},
		{path: "/nugget",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	url, cleanup := setupApi(t)
	defer cleanup()

	for _, testCase := range testCases {
		res, err := http.Get(url + testCase.path)

		if err != nil {
			t.Error(err)
		}

		defer res.Body.Close()

		if res.StatusCode != testCase.expectedStatusCode {
			t.Errorf("Expected %q, got %q", http.StatusText(testCase.expectedStatusCode), http.StatusText(res.StatusCode))
		}

		switch {
		case strings.Contains(res.Header.Get("Content-Type"), "text/html") || strings.Contains(res.Header.Get("Content-Type"), "text/plain"):
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}
			bodyParsed := string(body)

			if !strings.Contains(bodyParsed, testCase.expectedContent) {
				t.Errorf("Expected %q, got %q", testCase.expectedContent, bodyParsed)
			}
		default:
			t.Fatalf("Unsupported Content-Type: %q", res.Header.Get("Content-Type"))
		}
	}
}

func TestCreate(t *testing.T) {
	url, cleanup := setupApi(t)
	fmt.Println(url)
	defer cleanup()

	title := "The Little Peanut"
	content := "See how fast she runs!"

	t.Run("Add", func(t *testing.T) {
		var body bytes.Buffer
		item := struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}{
			Title:   title,
			Content: content,
		}

		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(url+"/snippet/create", "application/json", &body)

		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusCreated {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusCreated), http.StatusText(res.StatusCode))
		}
	})

	t.Run("Check Add", func(t *testing.T) {
		res, err := http.Get(url + "/snippet/view?id=1")
		if err != nil {
			t.Error(err)
		}

		var response snippetResponse

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
		}

		res.Body.Close()

		if response.Result.Title != title {
			t.Errorf("Expected title %q, got %q", title, response.Result.Title)
		}
		if response.Result.Content != content {
			t.Errorf("Expected title %q, got %q", content, response.Result.Content)
		}

	})
}
