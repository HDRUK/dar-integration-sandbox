package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// IHelper - BaseHelper interface
type IHelper interface {
	isAuthorized(ctx context.Context, token *string) bool
	getAccessToken(clientID string, clientSecret string, logger *zap.SugaredLogger) (string, error)
	httpRequest(method string, url string, body []byte) error
}

// BaseHelper - hold scoped logger, query
type BaseHelper struct {
	logger *zap.SugaredLogger
	query  IQuery
}

// NewBaseHelper - instantiate a new BaseHelper
func NewBaseHelper(query IQuery, logger *zap.SugaredLogger) *BaseHelper {
	return &BaseHelper{
		logger: logger,
		query:  query,
	}
}

// isAuthorized - validate bearer token
func (h *BaseHelper) isAuthorized(ctx context.Context, token *string) bool {
	filter := bson.M{"key": *token}

	// MongoDB find is wrapped in a separate service so that we can easily mock it in tests
	results := h.query.find(ctx, "tokens", filter)

	// In summary, if any accounts have this access token, then isAuthorized returns true
	if len(results) > 0 {
		return true
	}

	return false
}

type credentials struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// getAccessToken - retrieva a Gateway access token using the client credentials grant flow
func (h *BaseHelper) getAccessToken(clientID string, clientSecret string, logger *zap.SugaredLogger) (string, error) {
	credentials, _ := json.Marshal(&credentials{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	res, err := http.Post(os.Getenv("GATEWAY_BASE_URL")+"/oauth/token", "application/json", bytes.NewBuffer(credentials))
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer res.Body.Close()

	var tokenBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&tokenBody)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	accessToken := tokenBody["access_token"].(string)

	return accessToken, nil
}

// httpRequest -  simple wrapper for making a HTTP request to a target server, given a httpMethod
func (h *BaseHelper) httpRequest(method string, url string, body []byte) error {
	accessToken, err := h.getAccessToken(os.Getenv("GATEWAY_CLIENT_ID"), os.Getenv("GATEWAY_CLIENT_SECRET"), h.logger)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		return errors.New(method + " request to gateway-api received status code " + strconv.Itoa(res.StatusCode))
	}
	defer res.Body.Close()

	return nil
}
