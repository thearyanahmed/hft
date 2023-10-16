package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeleteService struct {
	mock.Mock
}

func (m *MockDeleteService) Delete(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockDeleteService) Exists(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func TestDeleteHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		urlParam     string
		expectedCode int
		mock         func(mockService *MockDeleteService)
	}{
		{
			name:         "Successful deletion",
			urlParam:     "TestConfig",
			expectedCode: http.StatusOK,
			mock: func(mockService *MockDeleteService) {
				mockService.On("Exists", "TestConfig").Return(true)
				mockService.On("Delete", "TestConfig").Return(nil)
			},
		},
		{
			name:         "Config not found",
			urlParam:     "NotFoundConfig",
			expectedCode: http.StatusNotFound,
			mock: func(mockService *MockDeleteService) {
				mockService.On("Exists", "NotFoundConfig").Return(false)
			},
		},
		{
			name:         "Error from DeleteService",
			urlParam:     "ErrorConfig",
			expectedCode: http.StatusUnprocessableEntity,
			mock: func(mockService *MockDeleteService) {
				mockService.On("Exists", "ErrorConfig").Return(true)
				mockService.On("Delete", "ErrorConfig").Return(errors.New("service error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockDeleteService)
			deleteHandlerInstance := NewDeleteHandler(mockService)

			tt.mock(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/configs/"+tt.urlParam, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Method(http.MethodDelete, "/configs/{name}", http.HandlerFunc(deleteHandlerInstance.ServeHTTP))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
