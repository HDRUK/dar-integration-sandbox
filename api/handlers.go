package api

import (
	"encoding/json"
	"net/http"

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
	TopicID      string      `json:"topicId" validate:"required"`
	MessageID    int64       `json:"messageId" validate:"required"`
	CreatedDate  string      `json:"createdDate"`
	QuestionBank interface{} `json:"questionBank"`
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
