package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type headerRecorder struct {
	http.ResponseWriter
	wroteHeader bool
}

func (r *headerRecorder) WriteHeader(code int) {
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *headerRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

func Recovery(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &headerRecorder{ResponseWriter: w}
			defer func() {
				if rec := recover(); rec != nil {
					rid := GetRequestID(r)
					stack := debug.Stack()
					logger.WithFields(logrus.Fields{
						"rid":   rid,
						"panic": rec,
					}).Errorf("panic recovered\n%s", string(stack))
					if !rec.wroteHeader {
						rec.Header().Set("Content-Type", "application/json")
						rec.WriteHeader(http.StatusInternalServerError)
						_ = json.NewEncoder(rec).Encode(map[string]string{"error": "internal error"})
					}
				}
			}()
			next.ServeHTTP(rec, r)
		})
	}
}
