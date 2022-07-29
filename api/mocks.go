package api

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// MockLogger - mock implementation of zap.Logger
func MockLogger(t *testing.T) (*zap.SugaredLogger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	slogger := logger.Sugar()
	return slogger, logs
}

// MockHelper - mock implementation of BaseHelper
type MockHelper struct{}

func (mh *MockHelper) isAuthorized(ctx context.Context, token *string) bool {
	if *token == "authorized" {
		return true
	}

	return false
}

func (mh *MockHelper) getAccessToken(clientID string, clientSecret string, logger *zap.SugaredLogger) (string, error) {
	return "", nil
}

// MockQuery - mock implementation of BaseQuery
type MockQuery struct{}

func (mq *MockQuery) find(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []bson.M {
	if filter["key"] == "authorized" {
		return []bson.M{{"name": "testAccount"}}
	}

	if filter["key"] == "timeout" {
		time.Sleep(6 * time.Second)
	}

	return []bson.M{}
}
