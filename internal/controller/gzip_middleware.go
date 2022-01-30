package controller

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipDecompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Type"), "gzip") {
			r.Header.Del("Content-Length")
			body, err := gzip.NewReader(r.Body)
			defer r.Body.Close()

			if err != nil {
				http.Error(w, "can't read gzipped data", http.StatusInternalServerError)
				return
			}
			r.Body = body

		}
		next.ServeHTTP(w, r)
	})
}
