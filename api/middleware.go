package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

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
func (b *BaseMiddleware) AuthorizeBearerToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var isAuthorized bool = false
		go func(ctx context.Context) {
			defer func() {
				recover()
			}()
			isAuthorized = b.service.isAuthorized(ctx, &authToken)

			cancel()
		}(ctx)

		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				b.Logger.Warn("Timeout waiting for MongoDB response")

				json.NewEncoder(w).Encode(
					&DefaultResponse{
						Success: false,
						Status:  "INTERNAL SERVER ERROR",
						Message: "Timeout waiting for MongoDB response",
					},
				)
			case context.Canceled:
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

				next(w, r)
			}
		}
	}
}
