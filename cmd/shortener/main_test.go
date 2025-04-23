package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	shortenerConfig "github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testBaseURL string = "http://my-awesome-shotener.com/"

func getTestShortener() *shortener.Shortener {
	conf := shortenerConfig.NewShortenerConfig()
	_ = conf.SetURLBase(testBaseURL)
	return shortener.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10), conf)
}

func TestGetShortLinkHandler(t *testing.T) {

	handler := handlers.GetShortLinkHandler(getTestShortener())

	srv := httptest.NewServer(handler)
	defer srv.Close()

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

			client := &http.Client{}
			req, _ := http.NewRequest(test.method, srv.URL, strings.NewReader(test.fullLink))
			resp, err := client.Do(req)
			if err != nil {
				assert.NoError(t, err)
			}
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

	shorter := getTestShortener()

	testCases := []struct {
		name         string
		shortLink    string
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
			shortLink:    testBaseURL + "1111",
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
			if test.shortLink == "" {
				short, err := shorter.GenerateShortLink(test.fullLink)
				assert.NoError(t, err)
				test.shortLink = short
				assert.NotEmpty(t, test.shortLink)
			}
			r := httptest.NewRequest(test.method, test.shortLink, nil)
			r.SetPathValue("id", strings.Replace(test.shortLink, testBaseURL, "", 1))

			handler := handlers.RedirectToFullLinkHandler(shorter)
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

func Test_getShortener(t *testing.T) {
	assert.IsType(t, &shortener.Shortener{}, getTestShortener())
}
