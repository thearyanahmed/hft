package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of FindService for testing
type MockFindService struct {
	mock.Mock
}

type testCase struct {
	name          string
	urlParam      string
	expectedCode  int
	expectedError error
	mock          func(mockService *MockFindService, tt testCase)
}

func (m *MockFindService) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	args := m.Called(options)
	return args.Get(0).([]schema.ConfigMap), args.Error(1)
}

func TestFindHandler_ServeHTTP(t *testing.T) {
	tests := []testCase{
		{
			name:         "Successful find",
			urlParam:     "TestConfig",
			expectedCode: http.StatusOK,
			mock: func(mockService *MockFindService, tt testCase) {
				mockService.On("Find", mock.Anything).Return([]schema.ConfigMap{{Name: tt.urlParam}}, tt.expectedError)
			},
		},
		{
			name:          "Error from FindService",
			urlParam:      "NotFoundConfig",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: errors.New("service error"),
			mock: func(mockService *MockFindService, tt testCase) {
				mockService.On("Find", mock.Anything).Return([]schema.ConfigMap{{Name: tt.urlParam}}, tt.expectedError)
			},
		},
		{
			name:          "No results found",
			urlParam:      "NonExistentConfig",
			expectedCode:  http.StatusNotFound,
			expectedError: nil,
			mock: func(mockService *MockFindService, tt testCase) {
				mockService.On("Find", mock.Anything).Return([]schema.ConfigMap{}, tt.expectedError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockFindService)
			findHandlerInstance := NewFindHandler(mockService) // Renamed variable

			tt.mock(mockService, tt)

			req := httptest.NewRequest(http.MethodGet, "/configs/"+tt.urlParam, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Method(http.MethodGet, "/configs/{name}", http.HandlerFunc(findHandlerInstance.ServeHTTP))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedError != nil {
				assert.Contains(t, w.Body.String(), tt.expectedError.Error())
			}

			mockService.AssertExpectations(t)
		})
	}
}
