package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

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

	// Capture the DAR application ID so that we can approve the DAR application
	applicationID := application["details"].(map[string]interface{})["dataRequestId"].(string)

	eg := new(errgroup.Group)

	eg.Go(func() error {
		// Get an access token from gateway-api using the service account credentials
		accessToken, err := GetAccessToken(os.Getenv("GATEWAY_CLIENT_ID"), os.Getenv("GATEWAY_CLIENT_SECRET"), h.Logger)
		if err != nil {
			return err
		}

		// Structure the JSON body that we need to send to approve a data access request
		messageToSend, _ := json.Marshal(map[string]string{
			"applicationStatus":            "approved",
			"applicationStatusDescription": "Approved automatically by the sandbox server!",
		})

		// Instantiate a new HTTP client
		client := &http.Client{}

		// Structure a PUT request to the gateway-api to approve an application
		req, err := http.NewRequest(http.MethodPut, os.Getenv("GATEWAY_BASE_URL")+"/api/v1/data-access-request/"+applicationID, bytes.NewBuffer(messageToSend))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		// Catch responses which are NOT 200 or 201
		if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
			return errors.New("DAR approval request to Gateway received status code " + strconv.Itoa(res.StatusCode))
		}
		defer res.Body.Close()

		return nil
	})

	// Using errgroup as akin to a try/catch, requires revision but allows us to send 500 in first instance
	if err := eg.Wait(); err != nil {
		h.Logger.Error(err)
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
			Message: "Data Access Request Submitted",
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

	validate := validator.New()
	if err := validate.Struct(message); err != nil {
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

	// The below is a repeat of the application handler logic, needs breaking out into re-usable components
	eg := new(errgroup.Group)

	eg.Go(func() error {
		accessToken, err := GetAccessToken(os.Getenv("GATEWAY_CLIENT_ID"), os.Getenv("GATEWAY_CLIENT_SECRET"), h.Logger)
		if err != nil {
			return err
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
			return err
		}

		if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
			return errors.New("Reply message request to Gateway received status code " + strconv.Itoa(res.StatusCode))
		}
		defer res.Body.Close()

		return nil
	})

	if err := eg.Wait(); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(
			&DefaultResponse{
				Success: false,
				Status:  "INTERNAL SERVER ERROR",
				Message: "Error sending reply to Gateway API",
			},
		)

		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		&DefaultResponse{
			Success: true,
			Status:  "OK",
			Message: "Data Access Request Submitted",
		},
	)

}
