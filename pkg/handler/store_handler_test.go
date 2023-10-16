package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	apiMiddleware "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/middleware"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of StoreService for testing
type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) Store(entity schema.ConfigMap) (schema.ConfigMap, error) {
	args := m.Called(entity)
	return args.Get(0).(schema.ConfigMap), args.Error(1)
}

func (m *MockStoreService) Exists(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func TestStoreHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      url.Values
		expectedCode     int
		expectedError    error
		expectedErrorMsg string
		mock             func(mockService *MockStoreService)
	}{
		{
			name: "Successful config creation",
			requestBody: url.Values{
				"name":     {"TestConfig"},
				"metadata": {`{"key0": "100", "key1":{"key2":"0g","key3":"1g"},"key4":{"key5":"4g","key6":"1g"},"key7":{"key8":"false","key9": "false","key10":"true"}}`},
			},
			expectedCode: http.StatusCreated,
			mock: func(mockService *MockStoreService) {
				mockService.On("Exists", "TestConfig").Return(false)
				mockService.On("Store", mock.AnythingOfType("schema.ConfigMap")).Return(schema.ConfigMap{
					Name: "TestConfig",
					Metadata: map[string]interface{}{
						"monitoring": map[string]interface{}{
							"enabled": "true",
						},
						"limits": map[string]interface{}{
							"cpu": map[string]interface{}{
								"enabled": "false",
								"value":   "300m",
							},
						},
					},
				}, nil)
			},
		},
		{
			name: "Config with same name already exists",
			requestBody: url.Values{
				"name":     {"TestConfig"},
				"metadata": {`{"key0": "100", "key1":{"key2":"0g","key3":"1g"},"key4":{"key5":"4g","key6":"1g"},"key7":{"key8":"false","key9": "false","key10":"true"}}`},
			},
			expectedCode:     http.StatusUnprocessableEntity,
			expectedErrorMsg: "config with name 'TestConfig' already exists",
			mock: func(mockService *MockStoreService) {
				mockService.On("Exists", "TestConfig").Return(true)
			},
		},
		{
			name:             "Validation error in request body",
			requestBody:      url.Values{},
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "validation failed",
			mock: func(mockService *MockStoreService) {
				// No mock setup for service calls, as this test checks for request body validation
			},
		},
		{
			name: "Error from StoreService",
			requestBody: url.Values{
				"name":     {"ErrorConfig"},
				"metadata": {`{"key0": "100", "key1":{"key2":"0g","key3":"1g"},"key4":{"key5":"4g","key6":"1g"},"key7":{"key8":"false","key9": "false","key10":"true"}}`},
			},
			expectedCode:     http.StatusUnprocessableEntity,
			expectedErrorMsg: "service error",
			mock: func(mockService *MockStoreService) {
				mockService.On("Exists", "ErrorConfig").Return(false)
				mockService.On("Store", mock.AnythingOfType("schema.ConfigMap")).Return(schema.ConfigMap{}, errors.New("service error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockStoreService)
			storeHandlerInstance := NewStoreHandler(mockService)

			tt.mock(mockService)

			reqBody := tt.requestBody.Encode()

			// Create a request with the specified request body
			req := httptest.NewRequest(http.MethodPost, "/configs", bytes.NewBufferString(reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			apiMiddleware.ValidateContentTypeMiddleware(storeHandlerInstance).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErrorMsg)
			}

			mockService.AssertExpectations(t)
		})
	}
}
