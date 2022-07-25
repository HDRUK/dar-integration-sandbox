package api

import (
	"context"
	"testing"

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

// MockService - mock implementation of BaseService
type MockService struct{}

func (ms *MockService) isAuthorized(ctx context.Context, token *string) bool {
	if *token == "authorized" {
		return true
	}

	return false
}

// MockQuery - mock implementation of BaseQuery
type MockQuery struct{}

func (mq *MockQuery) find(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []bson.M {
	if filter["key"] == "authorized" {
		return []bson.M{{"name": "testAccount"}}
	}

	return []bson.M{}
}
