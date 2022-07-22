package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// BaseMiddleware - hold db, logger
type BaseMiddleware struct {
	Logger  *zap.SugaredLogger
	service IService
}

// NewBaseMiddleware - instantiate a new BaseMiddleware
func NewBaseMiddleware(s IService, logger *zap.SugaredLogger) *BaseMiddleware {
	return &BaseMiddleware{
		Logger:  logger,
		service: s,
	}
}

// AuthorizeBearerToken - check the validity of a bearer token
func (b *BaseMiddleware) AuthorizeBearerToken(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]

		var isAuthorized bool = false
		isAuthorized = b.service.isAuthorized(&authToken)

		if !isAuthorized {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			b.Logger.Warn("Unauthorized access attempt")

			json.NewEncoder(w).Encode(
				&DefaultResponse{
					Success: false,
					Status:  "UNAUTHORIZED",
					Message: "You are not authorized to perform this request",
				},
			)

			return
		}

		f(w, r)
	}
}
