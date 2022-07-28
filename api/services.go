package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// IService - BaseService interface
type IService interface {
	isAuthorized(ctx context.Context, token *string) bool
}

// BaseService - hold scoped db
type BaseService struct {
	logger *zap.SugaredLogger
	query  IQuery
}

// NewBaseService - instantiate a new BaseService
func NewBaseService(query IQuery, logger *zap.SugaredLogger) *BaseService {
	return &BaseService{
		logger: logger,
		query:  query,
	}
}

// isAuthorized - validate bearer token
func (s *BaseService) isAuthorized(ctx context.Context, token *string) bool {
	filter := bson.M{"key": *token}

	// MongoDB find is wrapped in a separate service so that we can easily mock it in tests
	results := s.query.find(ctx, "tokens", filter)

	// In summary, if any accounts have this access token, then isAuthorized returns true
	if len(results) > 0 {
		return true
	}

	return false
}
