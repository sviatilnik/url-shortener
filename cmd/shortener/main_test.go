package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getTestShortener() *app.Shortener {
	return app.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10))
}

func TestGetShortLinkHandler(t *testing.T) {
	testCases := []struct {
		name         string
		fullLink     string
		expectedCode int
		method       string
	}{
		{
			name:         "#1",
			fullLink:     "http://google.com",
			expectedCode: http.StatusCreated,
			method:       http.MethodPost,
		},
		{
			name:         "#2",
			fullLink:     "",
			expectedCode: http.StatusBadRequest,
			method:       http.MethodPost,
		},
		{
			name:         "#3",
			fullLink:     " ",
			expectedCode: http.StatusBadRequest,
			method:       http.MethodPost,
		},
		{
			name:         "#4",
			fullLink:     "give me short link!",
			expectedCode: http.StatusBadRequest,
			method:       http.MethodPost,
		},
		{
			name:         "#5",
			fullLink:     "http://google.com",
			expectedCode: http.StatusMethodNotAllowed,
			method:       http.MethodGet,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.fullLink))

			handler := GetShortLinkHandler(getTestShortener())
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expectedCode, resp.StatusCode)

			if resp.StatusCode == http.StatusCreated {
				respBody, _ := io.ReadAll(resp.Body)
				assert.NotEmpty(t, respBody)
			}
		})
	}
}

func TestRedirectToFullLinkHandler(t *testing.T) {

	shortener := getTestShortener()

	testCases := []struct {
		name         string
		id           string
		fullLink     string
		expectedCode int
		method       string
	}{
		{
			name:         "#1",
			fullLink:     "http://google.com",
			expectedCode: http.StatusTemporaryRedirect,
			method:       http.MethodGet,
		},
		{
			name:         "#2",
			fullLink:     "http://google.com",
			id:           "1111",
			expectedCode: http.StatusBadRequest,
			method:       http.MethodGet,
		},
		{
			name:         "#5",
			fullLink:     "http://google.com",
			expectedCode: http.StatusMethodNotAllowed,
			method:       http.MethodPost,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if test.id == "" {
				short, err := shortener.GenerateShortLink(test.fullLink)
				assert.NoError(t, err)
				test.id = short
				assert.NotEmpty(t, test.id)
			}
			r := httptest.NewRequest(test.method, "/"+test.id, nil)
			r.SetPathValue("id", test.id)
			
			handler := RedirectToFullLinkHandler(shortener)
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expectedCode, resp.StatusCode)
			if resp.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, test.fullLink, resp.Header.Get("Location"))
			}
		})
	}
}
