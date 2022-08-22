package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// BaseMiddleware - hold db, logger
type BaseMiddleware struct {
	logger *zap.SugaredLogger
	helper IHelper
}

// NewBaseMiddleware - instantiate a new BaseMiddleware
func NewBaseMiddleware(helper IHelper, logger *zap.SugaredLogger) *BaseMiddleware {
	return &BaseMiddleware{
		logger: logger,
		helper: helper,
	}
}

// AuthorizeBearerToken - check the validity of a bearer token
func (b *BaseMiddleware) AuthorizeBearerToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context with a timeout so that the request doesn't hang is there is a MongoDB issue
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var isAuthorized bool = false

		// Initiate a go routine to determine if the auth token is valid
		go func(ctx context.Context) {
			defer func() {
				recover()
				cancel()
			}()
			authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]

			fmt.Println(authToken)

			isAuthorized = b.helper.isAuthorized(ctx, &authToken)

			cancel()
		}(ctx)

		// Determine the outcome of go routine, canceled our deadlineExceeded denotes error case
		// If all okay, pass request to handler
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				b.logger.Warn("Timeout waiting for MongoDB response")

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

					b.logger.Warn("Unauthorized access attempt")

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
