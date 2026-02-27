package repo

import (
	"context"
	"daybook/backend/models"

	"github.com/google/uuid"
)

type PageRepository interface {
	Create(ctx context.Context, page *models.Page) error
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Page, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Page, error)
	Update(ctx context.Context, page *models.Page) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
