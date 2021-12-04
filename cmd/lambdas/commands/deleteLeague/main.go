package main

import (
	"context"
	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/manage/leagues"
	"golang-serverless-cdk-template/internal/persistence"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func handlerFn(ctx context.Context, event events.APIGatewayProxyRequest) (*api.APIResponse, error) {
	db, err := persistence.OpenConnection(&persistence.OpenConnectionInput{
		DBHost:     os.Getenv("DB_HOST"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	})

	if err != nil {
		return nil, err
	}

	defer persistence.CloseConnection(db)

	service := leagues.NewService(db)

	deleteLeagueRequest, err := leagues.NewDeleteRequest(event)

	if err != nil {
		return nil, err
	}

	resp, err := service.DeleteLeague(deleteLeagueRequest)

	if err != nil {
		return nil, err
	}

	apiResp := api.NewAPIResponseBuilder().
		StatusCode(200).
		Header("Content-Type", "application/json").
		Data(resp).
		Build()

	return apiResp, nil
}

func main() {
	logrus.Info("Starting LambdaHandler")
	lambda.Start(
		api.LambdaResponder(handlerFn),
	)
}
