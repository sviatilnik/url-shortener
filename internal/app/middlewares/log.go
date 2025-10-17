package middlewares

import (
	"net/http"
	"time"

	"github.com/sviatilnik/url-shortener/internal/app/logger"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg, err := logger.NewLogger()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		start := time.Now()

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: 0,
				size:   0,
			},
		}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logg.Infow(
			"request",
			"uri", r.RequestURI,
			"method", r.Method,
			"request_Authorization", r.Header.Get("Authorization"),
			"request_Cookie", r.Header.Get("Cookie"),
			"status", lw.data.status,
			"duration", duration,
			"size", lw.data.size,
		)
	})
}

type responseData struct {
	size   int
	status int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.data.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.data.status = statusCode
}
