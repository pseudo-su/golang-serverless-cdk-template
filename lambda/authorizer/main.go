package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	logrus.Info("Running Authorizer.")

	if os.Getenv("USER_AUTH_ENABLED") != "true" {
		return generatePolicy("user", "Allow", event.MethodArn), nil
	}

	headerValue := event.AuthorizationToken
	logrus.Info("Validating token.")
	token, err := extractToken(headerValue)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token format. Expected format <Bearer <token>>")
	}
	if !isAuthenticated(token) {
		logrus.Info("Token Invalid.")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}

	logrus.Info("Token Valid.")

	return generatePolicy("user", "Allow", event.MethodArn), nil
}

func main() {
	lambda.Start(handleRequest)
}

func extractToken(t string) (string, error) {
	if len(t) == 0 {
		return "", fmt.Errorf("Invalid token")
	}
	splits := strings.Split(t, "Bearer")

	if len(splits) != 2 {
		return "", fmt.Errorf("Invalid token")
	}
	return strings.TrimSpace(splits[1]), nil
}

type OktaIntrospectResponse struct {
	Active               bool   `json:"active"`
	UserId               string `json:"uid"`
	Email                string `json:"username"`
	TokenIssueTimestamp  int64  `json:"iat"`
	TokenExpireTimestamp int64  `json:"exp"`
}

func isAuthenticated(t string) bool {
	endpoint := os.Getenv("USER_AUTH_ENDPOINT")
	cid := os.Getenv("USER_AUTH_CLIENT_ID")
	tokenTypeHint := os.Getenv("USER_AUTH_TOKEN_TYPE_HINT")

	url := endpoint + "?" + "client_id=" + cid + "&token_type_hint=" + tokenTypeHint + "&token=" + t

	res, err := http.Post(url, "application/x-www-form-urlencoded; charset=utf-8", nil)

	if err != nil {
		logrus.Error(err)
		return false
	}
	defer res.Body.Close()

	oktaIntrospectResponse := &OktaIntrospectResponse{}
	err = json.NewDecoder(res.Body).Decode(oktaIntrospectResponse)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return oktaIntrospectResponse.Active
}

func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{"*"},
				},
			},
		}
	}
	return authResponse
}
