// Модуль сжатия gzip для запросов/ответов.
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Получение заголовка
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Запись body
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Установка заголовка
func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Закрытие запроса
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Чтение body
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Закрытие запроса
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Мидлваря для сжатия запросов/ответов.
func GZip(acceptedContentTypes []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			contentType := r.Header.Get("Content-Type")
			suitableForGZip := false
			for _, v := range acceptedContentTypes {
				if strings.Contains(contentType, v) {
					suitableForGZip = true
					break
				}
			}

			if !suitableForGZip {
				next.ServeHTTP(w, r)
				return
			}

			ow := w

			// клиент принимает данные в gzip: сжимаем
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			// клиент прислал данные в gzip: распаковываем
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}

			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
