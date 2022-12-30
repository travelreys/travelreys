package utils

import "net/http"

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{w, http.StatusOK}
}

func (wrw *wrappedResponseWriter) WriteHeader(code int) {
	wrw.statusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}

type WrappedReponseWriterMiddleware struct{}

func NewWrappedReponseWriterMiddleware() *WrappedReponseWriterMiddleware {
	return &WrappedReponseWriterMiddleware{}
}

func (mw *WrappedReponseWriterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrw := NewWrappedResponseWriter(w)
		next.ServeHTTP(wrw, r)
	})
}
