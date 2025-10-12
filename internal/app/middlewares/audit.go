package middlewares

import (
	"net/http"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/audit"
	"github.com/sviatilnik/url-shortener/internal/app/models"
)

type AuditMiddleware struct {
	auditService *audit.AuditService
}

// ключ для хранения URL в контексте
type AuditContextKey string

const AuditURLKey AuditContextKey = "audit_url"

func NewAuditMiddleware(auditService *audit.AuditService) *AuditMiddleware {
	return &AuditMiddleware{
		auditService: auditService,
	}
}

func (m *AuditMiddleware) Audit(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем wrapper для ResponseWriter чтобы перехватить статус код
		wrapper := &responseWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Выполняем следующий handler
		nextHandler.ServeHTTP(wrapper, r)

		// Логируем событие только при успешном выполнении
		if wrapper.statusCode < 400 {
			m.logAuditEvent(r)
		}
	})
}

func (m *AuditMiddleware) logAuditEvent(r *http.Request) {

	userID, ok := r.Context().Value(models.ContextUserID).(string)
	if !ok {
		userID = ""
	}

	path := r.URL.Path
	method := r.Method

	switch {
	case method == "POST" && (path == "/" || path == "/api/shorten"):
		// Событие создания короткой ссылки
		url := m.extractURLFromRequest(r)
		if url != "" {
			m.auditService.LogShortenEvent(r.Context(), userID, url)
		}
	case method == "GET" && strings.HasPrefix(path, "/") && path != "/ping":
		// Событие перехода по короткой ссылке
		if url, ok := r.Context().Value(AuditURLKey).(string); ok && url != "" {
			m.auditService.LogFollowEvent(r.Context(), userID, url)
		}
	}
}

func (m *AuditMiddleware) extractURLFromRequest(r *http.Request) string {
	if url, ok := r.Context().Value(AuditURLKey).(string); ok {
		return url
	}
	return ""
}

// responseWrapper обертка для ResponseWriter для перехвата статус кода
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
