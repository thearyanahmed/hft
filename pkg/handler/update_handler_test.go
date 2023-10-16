package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of UpdateService for testing
type MockUpdateService struct {
	mock.Mock
}

func (m *MockUpdateService) Update(name string, entity schema.ConfigMap) (schema.ConfigMap, error) {
	args := m.Called(name, entity)
	return args.Get(0).(schema.ConfigMap), args.Error(1)
}

func (m *MockUpdateService) Exists(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name             string
		urlParam         string
		requestBody      url.Values
		expectedCode     int
		expectedError    error
		expectedErrorMsg string
		mock             func(mockService *MockUpdateService)
	}{
		{
			name:     "Successful config update",
			urlParam: "TestConfig",
			requestBody: url.Values{
				"name":     {"TestConfig"},
				"metadata": {`{"key0": "200", "key1":{"key2":"1g","key3":"0g"},"key4":{"key5":"3g","key6":"2g"},"key7":{"key8":"true","key9": "false","key10":"false"}}`},
			},
			mock: func(mockService *MockUpdateService) {
				mockService.On("Exists", "TestConfig").Return(true)
				mockService.On("Update", "TestConfig", mock.AnythingOfType("schema.ConfigMap")).Return(schema.ConfigMap{
					Name: "TestConfig",
					Metadata: map[string]interface{}{
						"key0": "200",
						"key1": map[string]interface{}{
							"key2": "1g",
							"key3": "0g",
						},
						"key4": map[string]interface{}{
							"key5": "3g",
							"key6": "2g",
						},
						"key7": map[string]interface{}{
							"key8":  "true",
							"key9":  "false",
							"key10": "false",
						},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "Config not found",
			urlParam: "NonExistentConfig",
			requestBody: url.Values{
				"name":     {"NonExistentConfig"},
				"metadata": {`{"key0": "200", "key1":{"key2":"1g","key3":"0g"},"key4":{"key5":"3g","key6":"2g"},"key7":{"key8":"true","key9": "false","key10":"false"}}`},
			},
			mock: func(mockService *MockUpdateService) {
				mockService.On("Exists", "NonExistentConfig").Return(false)
			},
			expectedCode:     http.StatusNotFound,
			expectedErrorMsg: "{\"message\":\"resource not found\"}\n",
		},
		{
			name:     "Error from UpdateService",
			urlParam: "ErrorConfig",
			requestBody: url.Values{
				"name":     {"ErrorConfig"},
				"metadata": {`{"key0": "200", "key1":{"key2":"1g","key3":"0g"},"key4":{"key5":"3g","key6":"2g"},"key7":{"key8":"true","key9": "false","key10":"false"}}`},
			},
			mock: func(mockService *MockUpdateService) {
				mockService.On("Exists", "ErrorConfig").Return(true)
				mockService.On("Update", "ErrorConfig", mock.AnythingOfType("schema.ConfigMap")).Return(schema.ConfigMap{}, errors.New("service error"))
			},
			expectedCode:     http.StatusUnprocessableEntity,
			expectedErrorMsg: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock UpdateService response
			mockService := new(MockUpdateService)
			updateHandlerInstance := NewUpdateHandler(mockService)

			tt.mock(mockService)

			reqBody := tt.requestBody.Encode()

			// Create a request with the specified request body
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/configs/%s", tt.requestBody.Get("name")), bytes.NewBufferString(reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			// Handle the request

			r := chi.NewRouter()
			r.Method(http.MethodPut, "/configs/{name}", http.HandlerFunc(updateHandlerInstance.ServeHTTP))
			r.ServeHTTP(w, req)
			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)
			// assert.Equal(t, tt.expectedCode, w.Body.String())

			// Reset the mock after each test case
			mockService.AssertExpectations(t)
		})
	}
}
