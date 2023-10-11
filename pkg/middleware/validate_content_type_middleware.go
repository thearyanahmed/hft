package middleware

import (
	"net/http"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
)

func ValidateContentTypeMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			presenter.ErrorResponse(w, r, presenter.ErrInvalidContentType())
			return
		}
		next.ServeHTTP(w, r.WithContext(r.Context()))
	}
	return http.HandlerFunc(fn)
}
