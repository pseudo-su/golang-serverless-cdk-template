package leagues

import (
	"encoding/json"
	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/persistence"
	"time"

	"github.com/aws/aws-lambda-go/events"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeleteDispatcher struct {
	leagueRepo *persistence.LeagueRepository
	validate   *validator.Validate
}

func NewDeleteDispatcher(db *gorm.DB) *DeleteDispatcher {
	return &DeleteDispatcher{
		leagueRepo: persistence.NewLeagueRepository(db),
		validate:   validator.New(),
	}
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type DeleteResponse struct {
	ID        uuid.UUID `json:"id"`
	DeletedAt time.Time `json:"deletedAt"`
}

func NewDeleteRequest(event events.APIGatewayProxyRequest) (*DeleteRequest, error) {
	var request DeleteRequest

	// Unmarshal the json, return 400 if error
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *DeleteDispatcher) DeleteLeague(req *DeleteRequest) (*DeleteResponse, error) {

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

	deletedLeague, err := s.leagueRepo.DeleteLeague(req.ID)

	if err != nil {
		return nil, err
	}

	return &DeleteResponse{
		ID:        deletedLeague.ID,
		DeletedAt: deletedLeague.DeletedAt.Time,
	}, nil
}
