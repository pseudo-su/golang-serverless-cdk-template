package leagues

import (
	"encoding/json"

	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/persistence"

	"github.com/aws/aws-lambda-go/events"
	validator "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CreateRequest struct {
	League *CreateValues `json:"league" validate:"required"`
}

type CreateResponse struct {
	League *SavedValues `json:"league"`
}

func NewCreateRequest(event events.APIGatewayProxyRequest) (*CreateRequest, error) {
	var request CreateRequest

	// Unmarshal the json, return 400 if error
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *Service) CreateLeague(req *CreateRequest) (*CreateResponse, error) {

	err := s.validate.Struct(req)

	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			logrus.Info(err)
			return nil, api.NewValidationError(err)
		} else {
			return nil, err
		}
	}

	// TODO: check user has permission to create league

	newLeague, err := s.leagueRepo.CreateLeague(&persistence.League{
		Name:        req.League.Name,
		Description: req.League.Name,
	})

	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		League: newSavedLeagueFromDB(newLeague),
	}, nil
}
