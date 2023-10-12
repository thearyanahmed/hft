package handler

import (
	"net/http"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/serializer"
)

type storeHandler struct {
	storeSvc StoreService
}

type StoreService interface {
	Store(entity schema.ConfigMap) (schema.ConfigMap, error)
	Exists(name string) bool
}

func NewStoreHandler(storeSvc StoreService) *storeHandler {
	return &storeHandler{storeSvc: storeSvc}
}

func (h *storeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.StoreConfigRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	data, err := eventRequest.ToConfigMapSchema()
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if h.storeSvc.Exists(data.Name) {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	createdConfig, err := h.storeSvc.Store(data)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusCreated, createdConfig)
}
