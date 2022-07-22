package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isAuthorized_Authorized(t *testing.T) {
	mockQuery := &MockQuery{}
	mockService := &BaseService{query: mockQuery}

	var token string = "authorized"
	authorized := mockService.isAuthorized(&token)

	assert.True(t, authorized)
}

func Test_isAuthorized_Unauthorized(t *testing.T) {
	mockQuery := &MockQuery{}
	mockService := &BaseService{query: mockQuery}

	var token string = "blah"
	authorized := mockService.isAuthorized(&token)

	assert.False(t, authorized)
}
