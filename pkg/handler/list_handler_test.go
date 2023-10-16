package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of ListService for testing
type MockListService struct {
	mock.Mock
}

func (m *MockListService) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	args := m.Called(options)
	return args.Get(0).([]schema.ConfigMap), args.Error(1)
}

func TestListHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		mock         func(mockService *MockListService)
	}{
		{
			name:         "Successful list retrieval",
			expectedCode: http.StatusOK,
			mock: func(mockService *MockListService) {
				mockService.On("Find", mock.Anything).Return([]schema.ConfigMap{{Name: "TestConfig"}}, nil)
			},
		},
		{
			name:         "Error from ListService",
			expectedCode: http.StatusUnprocessableEntity,
			mock: func(mockService *MockListService) {
				mockService.On("Find", mock.Anything).Return([]schema.ConfigMap{}, errors.New("service error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockListService)
			listHandlerInstance := NewListHandler(mockService)

			tt.mock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/configs", nil)
			w := httptest.NewRecorder()

			listHandlerInstance.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
