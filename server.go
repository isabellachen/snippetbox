package main

import "net/http"

func homeHandler(w http.ResponseWriter, res *http.Request) {
	if res.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	w.Write([]byte("hello from snippetbox"))
}

func snippetCreate(w http.ResponseWriter, res *http.Request) {
	if res.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func snippetView(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{name:"Isa"}`))
}

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.HandleFunc("/snippet/view", snippetView)
	return mux
}
