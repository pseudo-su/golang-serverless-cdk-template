package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeagueRepository struct {
	db *gorm.DB
}

//NewRepository create new repository
func NewLeagueRepository(db *gorm.DB) *LeagueRepository {
	return &LeagueRepository{
		db: db,
	}
}

func (r *LeagueRepository) CreateLeague(league *League) (*League, error) {
	result := r.db.Create(&league)

	if result.Error != nil {
		return nil, result.Error
	}

	return league, nil
}

func (r *LeagueRepository) GetLeagueByID(id uuid.UUID) (*League, error) {
	var league League

	result := r.db.First(&league, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &league, nil
}

func (r *LeagueRepository) UpdateLeague(league *League) (*League, error) {
	result := r.db.Model(&league).Updates(&league)

	if result.Error != nil {
		return nil, result.Error
	}

	return league, nil
}

func (r *LeagueRepository) DeleteLeague(id uuid.UUID) (*League, error) {
	league, err := r.GetLeagueByID(id)

	if err != nil {
		return nil, err
	}

	result := r.db.Delete(&league)

	if result.Error != nil {
		return nil, result.Error
	}

	return league, nil
}

type SearchLeaguesInput struct {
	*PaginationInput
}

type SearchLeaguesOutput struct {
	*PaginationOutput
	Leagues []League
}

func (r *LeagueRepository) SearchLeagues(input *SearchLeaguesInput) (*SearchLeaguesOutput, error) {
	leagues := []League{}
	output := &SearchLeaguesOutput{
		PaginationOutput: &PaginationOutput{},
	}

	result := r.db.Scopes(
		Paginate(leagues, r.db, input.PaginationInput, output.PaginationOutput),
	).Find(&leagues)

	if result.Error != nil {
		return nil, result.Error
	}

	output.Leagues = leagues

	return output, nil
}
