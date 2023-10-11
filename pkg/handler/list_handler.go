package handler

import (
	"net/http"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
)

type listHandler struct {
	listSvc listService
}

type listService interface{}

func NewListHandler(listSvc listService) *listHandler {
	return &listHandler{listSvc: listSvc}
}

func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	presenter.RenderJsonResponse(w, r, http.StatusOK, map[string]string{"sanity": "check"})
}
