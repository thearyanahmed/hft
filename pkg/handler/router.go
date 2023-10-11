package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apiMiddleware "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/middleware"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type ConfigManager interface {
	Store() error
	Find() ([]schema.ConfigMap, error)
	Update() error
	Delete() error
}

func NewRouter(svc ConfigManager) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				Get("/configs", NewListHandler(svc).ServeHTTP)

		})
	})

	return r
}
