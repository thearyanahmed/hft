package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apiMiddleware "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/middleware"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type ConfigManager interface {
	Store(entity schema.ConfigMap) (schema.ConfigMap, error)
	Find(options *schema.FilterOptions) ([]schema.ConfigMap, error)
	Update(entity schema.ConfigMap) error
	Delete(entity schema.ConfigMap) error
	Exists(name string) bool
}

func NewRouter(svc ConfigManager) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/configs", NewListHandler(svc).ServeHTTP)
	r.With(apiMiddleware.ValidateContentTypeMiddleware).Post("/configs", NewStoreHandler(svc).ServeHTTP)

	return r
}
