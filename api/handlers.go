package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

// DefaultResponse - default response struct
type DefaultResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BaseHandler - hold scoped db, logger
type BaseHandler struct {
	Logger  *zap.SugaredLogger
	service IService
}

// NewBaseHandler - instantiate a new BaseHandler
func NewBaseHandler(s IService, logger *zap.SugaredLogger) *BaseHandler {
	return &BaseHandler{
		Logger:  logger,
		service: s,
	}
}

// healthCheckHandler - test the API is up and running
func (h *BaseHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Server status check")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		&DefaultResponse{
			Success: true,
			Status:  "OK",
			Message: "The server is up and running",
		},
	)
}

// applicationHandler - handler for a DAR submission
func (h *BaseHandler) applicationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		&DefaultResponse{
			Success: true,
			Status:  "OK",
			Message: "Enquiry Message Submitted",
		},
	)
}

// expectedData - struct of expect first message schema
type firstMessageSchema struct {
	TopicID      string                 `json:"topicId" validate:"required"`
	MessageID    int64                  `json:"messageId" validate:"required"`
	CreatedDate  string                 `json:"createdDate"`
	QuestionBank map[string]interface{} `json:"questionBank"`
}

// firstMessageHandler - handler for a first message enquiry
func (h *BaseHandler) firstMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message firstMessageSchema
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(message); err != nil {
		w.Header().Set("Content-Type", "application/json")

		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			h.Logger.Warn(validationError)
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			&DefaultResponse{
				Success: false,
				Status:  "BAD REQUEST",
				Message: err.Error(),
			},
		)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		accessToken, err := GetAccessToken(os.Getenv("GATEWAY_CLIENT_ID"), os.Getenv("GATEWAY_CLIENT_SECRET"), h.Logger)
		if err != nil {
			cancel()
			wg.Done()
			return
		}

		messageToSend, _ := json.Marshal(map[string]interface{}{
			"messageType":        "message",
			"topic":              message.TopicID,
			"relatedObjectIds":   []string{message.QuestionBank["datasetsRequested"].([]interface{})[0].(map[string]interface{})["_id"].(string)},
			"messageDescription": "Hello from the sandbox server!",
		})

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPost, os.Getenv("GATEWAY_BASE_URL")+"/api/v1/messages", bytes.NewBuffer(messageToSend))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		res, err := client.Do(req)
		if err != nil {
			cancel()
			wg.Done()
			return
		}

		if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
			h.Logger.Error("Reply message request to Gateway received status code %d", res.StatusCode)
			cancel()
			wg.Done()
			return
		}

		defer res.Body.Close()
		wg.Done()
	}(ctx)
	wg.Wait()

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.Canceled:
			h.Logger.Error("Error sending reply to Gateway API")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(
				&DefaultResponse{
					Success: false,
					Status:  "INTERNAL SERVER ERROR",
					Message: "Error sending reply to Gateway API",
				},
			)

		}
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(
			&DefaultResponse{
				Success: true,
				Status:  "OK",
				Message: "Data Access Request Submitted",
			},
		)
	}
}
