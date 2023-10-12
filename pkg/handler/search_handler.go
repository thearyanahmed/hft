package handler

import (
	"net/http"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/presenter"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type searchHandler struct {
	searchSvc SearchService
}

type SearchService interface {
	Find(options *schema.FilterOptions) ([]schema.ConfigMap, error)
}

func NewSearchHandler(searchSvc SearchService) *searchHandler {
	return &searchHandler{searchSvc: searchSvc}
}

func (h *searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	queryMap := make(map[string]string)
	for key, val := range queryParams {
		queryMap[key] = val[0]
	}

	options := &schema.FilterOptions{
		SelectAllIfConditionsAreEmpty: false,
		Limit:                         100,
		Conditions:                    queryMap,
	}

	configs, err := h.searchSvc.Find(options)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, configs)
}
