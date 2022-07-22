package api

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// IService - BaseService interface
type IService interface {
	isAuthorized(str *string) bool
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
func (s *BaseService) isAuthorized(token *string) bool {
	filter := bson.M{"key": *token}

	results := s.query.find("tokens", filter)

	if len(results) > 0 {
		return true
	}

	return false
}
