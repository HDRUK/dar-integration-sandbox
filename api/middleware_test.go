package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func Test_AuthorizeBearerToken_Unauthorized(t *testing.T) {
	var path string = "/test-url"
	req := httptest.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()

	var token string = "Bearer unauthorized"
	req.Header.Set("Authorization", token)

	log, logs := MockLogger(t)

	var mockHelper = &MockHelper{}

	mockBaseMiddleware := &BaseMiddleware{Logger: log, helper: mockHelper}

	mockBaseMiddleware.AuthorizeBearerToken(mockHandler)(res, req)

	expectedJSONResponse := DefaultResponse{
		Success: false,
		Status:  "UNAUTHORIZED",
		Message: "You are not authorized to perform this request",
	}

	var resBody DefaultResponse
	json.NewDecoder(res.Body).Decode(&resBody)

	assert.Equal(t, req.URL.Path, path)
	assert.Equal(t, req.Header.Get("Authorization"), token)
	assert.Equal(t, res.Code, http.StatusUnauthorized)
	assert.Equal(t, resBody, expectedJSONResponse)
	assert.Equal(t, logs.Len(), 1)
	assert.Equal(t, logs.All()[0].Message, "Unauthorized access attempt")
}

func Test_AuthorizeBearerToken_Authorized(t *testing.T) {
	var path string = "/test-url"
	req := httptest.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()

	var token string = "Bearer authorized"
	req.Header.Set("Authorization", token)

	var mockHelper = &MockHelper{}

	mockBaseMiddleware := &BaseMiddleware{helper: mockHelper}

	mockBaseMiddleware.AuthorizeBearerToken(mockHandler)(res, req)

	assert.Equal(t, req.URL.Path, path)
	assert.Equal(t, req.Header.Get("Authorization"), token)
	assert.Equal(t, res.Code, http.StatusOK)
}
