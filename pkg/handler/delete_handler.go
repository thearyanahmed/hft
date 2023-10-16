package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
)

type deleteHandler struct {
	deleteSvc DeleteService
}

type DeleteService interface {
	Delete(name string) error
	Exists(name string) bool
}

func NewDeleteHandler(deleteSvc DeleteService) *deleteHandler {
	return &deleteHandler{deleteSvc: deleteSvc}
}

func (h *deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if name == "" {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(fmt.Errorf("a required parameter 'name' is missing")))
		return
	}

	exists := h.deleteSvc.Exists(name)

	if !exists {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	err := h.deleteSvc.Delete(name)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	res := presenter.Response{
		Message:        "config deleted successfully",
		HttpStatusCode: http.StatusOK,
	}

	presenter.RenderJsonResponse(w, r, res.HttpStatusCode, res)
}
