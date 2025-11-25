package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/config"
)

// TrustedSubnetMiddleware проверяет, что IP-адрес клиента из заголовка X-Real-IP
// входит в доверенную подсеть (CIDR). Если trusted_subnet пустое или IP не входит в подсеть,
// возвращает статус 403 Forbidden.
func TrustedSubnetMiddleware(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если trusted_subnet пустое, доступ запрещён
			if strings.TrimSpace(config.TrustedSubnet) == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Получаем IP-адрес из заголовка X-Real-IP
			clientIP := r.Header.Get("X-Real-IP")
			if clientIP == "" {
				// Если заголовок отсутствует, используем RemoteAddr
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					clientIP = r.RemoteAddr
				} else {
					clientIP = ip
				}
			}

			// Парсим IP-адрес
			ip := net.ParseIP(clientIP)
			if ip == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Парсим CIDR подсеть
			_, ipNet, err := net.ParseCIDR(config.TrustedSubnet)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Проверяем, входит ли IP в подсеть
			if !ipNet.Contains(ip) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// IP входит в доверенную подсеть, продолжаем обработку
			next.ServeHTTP(w, r)
		})
	}
}

