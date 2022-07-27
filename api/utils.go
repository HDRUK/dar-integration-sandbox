package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type credentials struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// GetAccessToken - retrieva a Gateway access token using the client credentials grant flow
func GetAccessToken(clientID string, clientSecret string, logger *zap.SugaredLogger) (string, error) {
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
