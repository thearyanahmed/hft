package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/serializer"
)

type updateHandler struct {
	updateSvc UpdateService
}

type UpdateService interface {
	Update(name string, entity schema.ConfigMap) (schema.ConfigMap, error)
	Exists(name string) bool
}

func NewUpdateHandler(updateSvc UpdateService) *updateHandler {
	return &updateHandler{updateSvc: updateSvc}
}

func (h *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if name == "" {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(fmt.Errorf("a required parameter 'name' is missing")))
		return
	}

	updateReq := &serializer.UpdateConfigRequest{}

	if formErrors := serializer.ValidatePostForm(r, updateReq); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	data, err := updateReq.ToConfigMapSchema()
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if !h.updateSvc.Exists(name) {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	updatedConfig, err := h.updateSvc.Update(name, data)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, updatedConfig)
}
