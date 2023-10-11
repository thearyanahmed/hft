package handler

import (
	"net/http"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type listHandler struct {
	listSvc ListService
}

type ListService interface {
	Find() ([]schema.ConfigMap, error)
}

func NewListHandler(listSvc ListService) *listHandler {
	return &listHandler{listSvc: listSvc}
}

func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	configList, err := h.listSvc.Find()

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, configList)
}
