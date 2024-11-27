package middleware

import (
	"net"
	"net/http"
)

// Проверяет доступ по маске подсети из заголовка X-Real-IP
func TrustedSubnet(ts string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.Header.Get("X-Real-IP")

			// Настройка открыта для всех - можно
			if ts == "" {
				next.ServeHTTP(w, r)
			}

			// Настройка закрыта по маске, но X-Real-IP пустой - нельзя
			if clientIP == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Ошибка при парсинге - нельзя
			_, cidr, err := net.ParseCIDR(ts)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// IP не входит в доверенную сеть - нельзя
			if !cidr.Contains(net.ParseIP(clientIP)) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// МОЖНА!
			next.ServeHTTP(w, r)
		})
	}
}
