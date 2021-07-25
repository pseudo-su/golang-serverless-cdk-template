package leagues

import (
	"golang-serverless-cdk-template/internal/persistence"

	validator "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	CreateLeague(*CreateRequest) (*CreateResponse, error)
	GetLeagueByID(*GetByIDRequest) (*GetByIDResponse, error)
	SearchLeagues(*SearchRequest) (*SearchResponse, error)
	UpdateLeague(*UpdateRequest) (*UpdateResponse, error)
	DeleteLeague(*DeleteRequest) (*DeleteResponse, error)
}

type Service struct {
	leagueRepo *persistence.LeagueRepository
	validate   *validator.Validate
}

var _ ServiceInterface = &Service{}

func NewService(db *gorm.DB) ServiceInterface {
	return &Service{
		leagueRepo: persistence.NewLeagueRepository(db),
		validate:   validator.New(),
	}
}
