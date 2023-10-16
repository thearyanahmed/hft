package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	args := m.Called(options)
	return args.Get(0).([]schema.ConfigMap), args.Error(1)
}

func TestSearchHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  url.Values
		expectedCode int
		mock         func(mockService *MockSearchService)
	}{
		{
			name: "Successful search with parameters",
			queryParams: url.Values{
				"name": {"TestConfig"},
			},
			mock: func(mockService *MockSearchService) {
				mockService.On("Find", mock.AnythingOfType("*schema.FilterOptions")).Return([]schema.ConfigMap{{Name: "TestConfig"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Successful search without parameters",
			queryParams: url.Values{},
			mock: func(mockService *MockSearchService) {
				mockService.On("Find", mock.AnythingOfType("*schema.FilterOptions")).Return([]schema.ConfigMap{{Name: "TestConfig"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Error from SearchService",
			queryParams: url.Values{
				"name": {"ErrorConfig"},
			},
			mock: func(mockService *MockSearchService) {
				mockService.On("Find", mock.AnythingOfType("*schema.FilterOptions")).Return([]schema.ConfigMap{}, errors.New("service error"))
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			searchHandlerInstance := NewSearchHandler(mockService)

			tt.mock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/configs?"+tt.queryParams.Encode(), nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Method(http.MethodGet, "/configs", http.HandlerFunc(searchHandlerInstance.ServeHTTP))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
