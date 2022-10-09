package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()
	testServer := httptest.NewServer(newMux())

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
			expectedContent:    "hello from snippetbox",
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
		case strings.Contains(res.Header.Get("Content-Type"), "text/plain"):
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
