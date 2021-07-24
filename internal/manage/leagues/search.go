package leagues

import (
	"golang-serverless-cdk-template/internal/api"
	"golang-serverless-cdk-template/internal/persistence"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	validator "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SearchDispatcher struct {
	leagueRepo *persistence.LeagueRepository
	validate   *validator.Validate
}

func NewSearchDispatcher(db *gorm.DB) *SearchDispatcher {
	return &SearchDispatcher{
		leagueRepo: persistence.NewLeagueRepository(db),
		validate:   validator.New(),
	}
}

type SearchRequest struct {
	Page  uint `json:"page" validate:"required,gte=1"`
	Limit uint `json:"limit" validate:"required,gte=1,lte=50"`
}

type PaginationValues struct {
	Page       uint `json:"page"`
	Limit      uint `json:"limit"`
	TotalItems uint `json:"totalItems"`
	TotalPages uint `json:"totalPages"`
}

type SearchResponse struct {
	Pagination *PaginationValues `json:"pagination" validate:"required"`
	Leagues    []*SavedValues    `json:"leagues"`
}

func NewSearchRequest(event events.APIGatewayProxyRequest) (*SearchRequest, error) {
	var request SearchRequest
	page, err := strconv.Atoi(event.QueryStringParameters["page"])

	if err != nil {
		return nil, err
	}

	request.Page = uint(page)

	limit, err := strconv.Atoi(event.QueryStringParameters["limit"])

	if err != nil {
		return nil, err
	}

	request.Limit = uint(limit)

	return &request, nil
}

func (s *SearchDispatcher) SearchLeagues(req *SearchRequest) (*SearchResponse, error) {

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

	output, err := s.leagueRepo.SearchLeagues(&persistence.SearchLeaguesInput{
		PaginationInput: &persistence.PaginationInput{
			Page:  req.Page,
			Limit: req.Limit,
		},
	})

	if err != nil {
		return nil, err
	}
	return newSearchResponse(output), nil
}

func newSearchResponse(searchResult *persistence.SearchLeaguesOutput) *SearchResponse {
	pagination := &PaginationValues{
		Page:       searchResult.Page,
		Limit:      searchResult.Limit,
		TotalItems: searchResult.TotalItems,
		TotalPages: searchResult.TotalPages,
	}
	leagues := []*SavedValues{}

	for _, dbLeague := range searchResult.Leagues {
		leagues = append(leagues, newSavedLeagueFromDB(&dbLeague))
	}

	return &SearchResponse{
		Pagination: pagination,
		Leagues:    leagues,
	}
}
