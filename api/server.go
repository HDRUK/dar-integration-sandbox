package api

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Router - handler for HTTP requests
type Router struct {
	mux *mux.Router
}

//NewServer - initiate a new http server
func NewServer(s *mux.Router, db *mongo.Database, logger *zap.SugaredLogger) *Router {
	r := Router{s}
	r.handleRoutes(db, logger)

	return &r
}

// handleRoutes - bind routes to the handler
func (r *Router) handleRoutes(db *mongo.Database, logger *zap.SugaredLogger) {
	baseQuery := NewBaseQuery(db)
	baseUtility := NewBaseHelper(baseQuery, logger)
	baseHandler := NewBaseHandler(baseUtility, logger)
	baseMiddleware := NewBaseMiddleware(baseUtility, logger)

	// Health check
	r.mux.HandleFunc("/status", baseHandler.healthCheckHandler).Methods("GET", "OPTIONS")

	// Business
	r.mux.HandleFunc("/first-message", baseMiddleware.AuthorizeBearerToken(baseHandler.firstMessageHandler)).Methods("POST", "OPTIONS")
	r.mux.HandleFunc("/application", baseMiddleware.AuthorizeBearerToken(baseHandler.applicationHandler)).Methods("POST", "OPTIONS")

	// Testing
	r.mux.HandleFunc("/error", baseMiddleware.AuthorizeBearerToken(baseHandler.errorHandler)).Methods("POST", "OPTIONS")

}
