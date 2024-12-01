package middleware

import (
	"net/http"

	"github.com/mikesvis/short/internal/subnet"
)

// Проверяет доступ по маске подсети из заголовка X-Real-IP
func TrustedSubnet(ts string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			clientIP := r.Header.Get("X-Real-IP")
			if !subnet.ValidateSubnet(clientIP, ts) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
