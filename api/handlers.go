package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// DefaultResponse - default response struct
type DefaultResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BaseHandler - hold scoped db, logger
type BaseHandler struct {
	logger *zap.SugaredLogger
	helper IHelper
}

// NewBaseHandler - instantiate a new BaseHandler
func NewBaseHandler(helper IHelper, logger *zap.SugaredLogger) *BaseHandler {
	return &BaseHandler{
		logger: logger,
		helper: helper,
	}
}

// healthCheckHandler - test the API is up and running
func (h *BaseHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Server status check")

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

// applicationHandler - handler for a DAR application submission
func (h *BaseHandler) applicationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the incoming request body
	var application map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("DAR application received: ", r.Body)

	eg := new(errgroup.Group)

	eg.Go(func() error {
		// Capture the DAR application ID so that we can approve the DAR application
		applicationID := application["dataRequestId"].(string)

		// Structure the JSON body that we need to send to approve a data access request
		messageToSend, _ := json.Marshal(map[string]string{
			"applicationStatus":            "approved",
			"applicationStatusDescription": "Approved automatically by the sandbox server!",
		})

		// Make a PUT request to the Gateway to automatically approve the data access request
		err := h.helper.httpRequest(http.MethodPut, os.Getenv("GATEWAY_BASE_URL")+"/api/v1/data-access-request/"+applicationID, messageToSend)
		if err != nil {
			return err
		}

		return nil
	})

	// Catch any errors in the above goroutine
	if err := eg.Wait(); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(
			&DefaultResponse{
				Success: false,
				Status:  "INTERNAL SERVER ERROR",
				Message: "Error updating application status on Gateway",
			},
		)

		return
	}
	// If errors == <nil> send OK response
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		&DefaultResponse{
			Success: true,
			Status:  "OK",
			Message: "Data Access Request submitted",
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
	w.Header().Set("Content-Type", "application/json")

	var message firstMessageSchema
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("First message enquiry received: ", message)

	validate := validator.New()
	if err := validate.Struct(message); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			h.logger.Warn(validationError)
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

	// The below is a repeat of the application handler logic, needs breaking out into re-usable components
	eg := new(errgroup.Group)

	eg.Go(func() error {
		messageToSend, _ := json.Marshal(map[string]interface{}{
			"messageType":        "message",
			"topic":              message.TopicID,
			"relatedObjectIds":   []string{message.QuestionBank["datasetsRequested"].([]interface{})[0].(map[string]interface{})["_id"].(string)},
			"messageDescription": "Hello from the sandbox server!",
		})

		err := h.helper.httpRequest(http.MethodPost, os.Getenv("GATEWAY_BASE_URL")+"/api/v1/messages", messageToSend)
		if err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(
			&DefaultResponse{
				Success: false,
				Status:  "INTERNAL SERVER ERROR",
				Message: "Error sending message reply to Gateway API",
			},
		)

		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		&DefaultResponse{
			Success: true,
			Status:  "OK",
			Message: "First message enquiry submitted",
		},
	)

}
