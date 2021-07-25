package leagues

import (
	"encoding/json"
	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/persistence"

	"github.com/aws/aws-lambda-go/events"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UpdateRequest struct {
	ID     uuid.UUID     `json:"id" validate:"required"`
	League *UpdateValues `json:"league" validate:"required"`
}

type UpdateResponse struct {
	League *SavedValues `json:"league"`
}

func NewUpdateRequest(event events.APIGatewayProxyRequest) (*UpdateRequest, error) {
	var request UpdateRequest

	// Unmarshal the json, return 400 if error
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *Service) UpdateLeague(req *UpdateRequest) (*UpdateResponse, error) {

	err := s.validate.Struct(req)

	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			logrus.Info(err)
			return nil, api.NewValidationError(err)
		} else {
			return nil, err
		}
	}

	// TODO: check user has permission to update league

	leagueToUpdate := &persistence.League{
		ID: req.ID,
	}
	if req.League.Name != nil {
		leagueToUpdate.Name = *req.League.Name
	}
	if req.League.Description != nil {
		leagueToUpdate.Description = *req.League.Description
	}

	updatedLeague, err := s.leagueRepo.UpdateLeague(leagueToUpdate)

	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		League: newSavedLeagueFromDB(updatedLeague),
	}, nil
}
