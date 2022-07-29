package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isAuthorized_Authorized(t *testing.T) {
	mockQuery := &MockQuery{}
	mockHelper := &BaseHelper{query: mockQuery}

	var token string = "authorized"
	authorized := mockHelper.isAuthorized(context.TODO(), &token)

	assert.True(t, authorized)
}

func Test_isAuthorized_Unauthorized(t *testing.T) {
	mockQuery := &MockQuery{}
	mockHelper := &BaseHelper{query: mockQuery}

	var token string = "blah"
	authorized := mockHelper.isAuthorized(context.TODO(), &token)

	assert.False(t, authorized)
}
