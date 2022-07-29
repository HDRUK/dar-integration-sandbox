package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/status", nil)
	res := httptest.NewRecorder()

	log, logs := MockLogger(t)

	mockBaseHandler := &BaseHandler{logger: log}

	mockBaseHandler.healthCheckHandler(res, req)

	expectedJSONResponse := DefaultResponse{
		Success: true,
		Status:  "OK",
		Message: "The server is up and running",
	}

	var resBody DefaultResponse
	json.NewDecoder(res.Body).Decode(&resBody)

	assert.Equal(t, req.URL.Path, "/status")
	assert.Equal(t, res.Code, http.StatusOK)
	assert.Equal(t, resBody, expectedJSONResponse)
	assert.Equal(t, logs.Len(), 1)
	assert.Equal(t, logs.All()[0].Message, "Server status check")
}
