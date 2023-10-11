package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apiMiddleware "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	demoSvcStruct := struct{}{}

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				Get("/configs", NewListHandler(demoSvcStruct).ServeHTTP)

		})
	})

	return r
}
