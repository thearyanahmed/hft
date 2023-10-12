package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type findHandler struct {
	findSvc FindService
}

type FindService interface {
	Find(options *schema.FilterOptions) ([]schema.ConfigMap, error)
}

func NewFindHandler(findSvc FindService) *findHandler {
	return &findHandler{findSvc: findSvc}
}

func (h *findHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if name == "" {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(fmt.Errorf("a required parameter 'name' is missing")))
		return
	}

	options := &schema.FilterOptions{
		SelectAllIfConditionsAreEmpty: false,
		Limit:                         1,
		Conditions: map[string]string{
			"name": name,
		},
	}

	configs, err := h.findSvc.Find(options)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if len(configs) == 0 {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, configs[0])
}
