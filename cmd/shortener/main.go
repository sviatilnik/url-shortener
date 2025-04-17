package main

import (
	"github.com/sviatilnik/url-shortener/internal/app"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"io"
	"net/http"
)

func GetShortLinkHandler(shortener *app.Shortener) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		url, err := io.ReadAll(r.Body)
		if err != nil || !util.IsURL(string(url)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortID, err := shortener.GenerateShortLink(string(url))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("http://" + r.Host + "/" + shortID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func RedirectToFullLinkHandler(shortener *app.Shortener) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")

		fullLink, err := shortener.GetLinkByID(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
	})
}

func main() {
	shortener := getShortener()
	mux := http.NewServeMux()
	mux.Handle("/", GetShortLinkHandler(shortener))
	mux.Handle("/{id}", RedirectToFullLinkHandler(shortener))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func getShortener() *app.Shortener {
	return app.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10))
}
