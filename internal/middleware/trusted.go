package middleware

import (
	"net"
	"net/http"
)

func TrustedSubnet(ts string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.Header.Get("X-Real-IP")

			if ts == "" || clientIP == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, cidr, err := net.ParseCIDR(ts)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !cidr.Contains(net.ParseIP(clientIP)) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
