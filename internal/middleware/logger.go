// Модуль логирования запросов в приложении.
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Логирование запросов и времени обработки.
func RequestResponseLogger(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			uri := r.RequestURI
			method := r.Method
			next.ServeHTTP(&lw, r)
			duration := time.Since(start)
			log.Infow(
				"Incoming request",
				"uri", uri,
				"method", method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		}
		return http.HandlerFunc(fn)
	}
}

// Запись размера в ответ
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// Запись заголовка
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
