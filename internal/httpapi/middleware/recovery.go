package middleware

import (
	"MiniJira/internal/httpapi"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func Recovery(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					rid := GetRequestID(r)
					stack := debug.Stack()
					logger.WithFields(logrus.Fields{
						"rid":   rid,
						"panic": rec,
					}).Errorf("panic recovered\n%s", string(stack))
					httpapi.WriteError(w, http.StatusInternalServerError, "internal error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
