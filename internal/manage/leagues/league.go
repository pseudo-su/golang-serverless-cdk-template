package leagues

import (
	"time"

	"golang-serverless-cdk-template/internal/persistence"

	"github.com/google/uuid"
)

type CreateValues struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateValues struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type SavedValues struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description" validate:"required"`
}

func newSavedLeagueFromDB(dbLeague *persistence.League) *SavedValues {
	return &SavedValues{
		ID:          dbLeague.ID,
		CreatedAt:   dbLeague.CreatedAt,
		UpdatedAt:   dbLeague.UpdatedAt,
		Name:        dbLeague.Name,
		Description: dbLeague.Description,
	}
}
