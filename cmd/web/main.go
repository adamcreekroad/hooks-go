package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/adamcreekroad/hooks-go/plex"
)

func plex_hook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/plex" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}

	if err := r.ParseMultipartForm(r.ContentLength); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	payload := r.FormValue("payload")

	plex.ProcessHook(payload)
}

func main() {
	http.HandleFunc("/plex", plex_hook)

	addr := fmt.Sprintf("%s:%s", config.Binding(), config.Port())

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
