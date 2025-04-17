package main

import (
	"github.com/sviatilnik/url-shortener/internal/app"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"io"
	"net/http"
)

func GetShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	url, err := io.ReadAll(r.Body)
	if err != nil || !util.IsUrl(string(url)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortId, err := app.GenerateShortLink(string(url))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://" + r.Host + "/" + shortId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetFullLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	fullLink, err := app.GetLinkById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", GetShortLinkHandler)
	mux.HandleFunc("/{id}", GetFullLinkHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
