package api

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

type APIHandlerFn func(ctx context.Context, event events.APIGatewayProxyRequest) (*APIResponse, error)
type LambdaHandlerFn func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func LambdaResponder(handlerFn APIHandlerFn) LambdaHandlerFn {
	return func(
		ctx context.Context,
		event events.APIGatewayProxyRequest,
	) (*events.APIGatewayProxyResponse, error) {
		resp, err := handlerFn(ctx, event)

		if err != nil {
			logrus.Error(err)
			return sendAPIResponseError(err)
		}

		return sendAPIResponse(resp)
	}
}

func sendAPIResponseError(err error) (*events.APIGatewayProxyResponse, error) {
	builder := NewAPIResponseBuilder().
		Header("Content-Type", "application/json")

	if errorer, ok := err.(Errorer); ok {
		builder.Error(errorer)
	} else {
		builder.Error(&BaseError{
			originalError: err,
			status:        500,
			code:          "InternalServiceError",
			title:         "Internal service error",
		})
	}

	return sendAPIResponse(builder.Build())
}

func sendAPIResponse(resp *APIResponse) (*events.APIGatewayProxyResponse, error) {

	responseBody, err := json.Marshal(resp.Body)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       string(responseBody),
	}, nil
}
