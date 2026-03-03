package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *StatusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *StatusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	return r.ResponseWriter.Write(b)
}

func Logging(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := GetRequestID(r)
			start := time.Now()

			rec := &StatusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			dur := time.Since(start)

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.WithFields(logrus.Fields{
				"rid":      rid,
				"status":   status,
				"duration": dur.Round(time.Millisecond),
				"method":   r.Method,
				"path":     r.URL.RequestURI(),
			}).Info("request")
		})
	}
}
