package leagues

import (
	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/persistence"

	"github.com/aws/aws-lambda-go/events"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetByIDDispatcher struct {
	leagueRepo *persistence.LeagueRepository
	validate   *validator.Validate
}

func NewGetByIDDispatcher(db *gorm.DB) *GetByIDDispatcher {
	return &GetByIDDispatcher{
		leagueRepo: persistence.NewLeagueRepository(db),
		validate:   validator.New(),
	}
}

type GetByIDRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type GetByIDResponse struct {
	*SavedValues
}

func NewGetByIDRequest(event events.APIGatewayProxyRequest) (*GetByIDRequest, error) {
	var request GetByIDRequest
	uuid, err := uuid.Parse(event.PathParameters["leagueId"])

	if err != nil {
		return nil, err
	}

	request.ID = uuid

	return &request, nil
}

func (s *GetByIDDispatcher) GetLeagueByID(req *GetByIDRequest) (*GetByIDResponse, error) {

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

	league, err := s.leagueRepo.GetLeagueByID(req.ID)

	if err != nil {
		return nil, err
	}

	return &GetByIDResponse{
		newSavedLeagueFromDB(league),
	}, nil
}
