package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const HeaderRequestID = "X-Request-Id"

type ctxKeyRequestID struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = NewRequestID()
		}

		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, rid)
		r.Header.Set(HeaderRequestID, rid)
		w.Header().Set(HeaderRequestID, rid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewRequestID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}

	return hex.EncodeToString(buf)
}

func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(ctxKeyRequestID{}).(string); ok && id != "" {
		return id
	}
	return r.Header.Get(HeaderRequestID)
}
