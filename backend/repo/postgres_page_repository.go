package repo

import (
	"context"
	"daybook/backend/db"
	"daybook/backend/models"

	"github.com/google/uuid"
)

type PostgresPageRepository struct{}

func NewPostgresPageRepository() *PostgresPageRepository {
	return &PostgresPageRepository{}
}

func (r *PostgresPageRepository) Create(ctx context.Context, page *models.Page) error {
	query := `
		INSERT INTO pages (id, user_id, title, content)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.Pool.Exec(ctx, query,
		page.ID,
		page.UserID,
		page.Title,
		page.Content,
	)

	return err
}

func (r *PostgresPageRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Page, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM pages
		WHERE user_id = $1
		AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`

	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pages := []models.Page{}

	for rows.Next() {
		var p models.Page
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (r *PostgresPageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM pages
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var p models.Page

	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.UserID,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PostgresPageRepository) Update(ctx context.Context, page *models.Page) error {
	query := `
		UPDATE pages
		SET title = $1,
		    content = $2
		WHERE id = $3
		AND deleted_at IS NULL
	`

	_, err := db.Pool.Exec(ctx, query,
		page.Title,
		page.Content,
		page.ID,
	)

	return err
}

func (r *PostgresPageRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE pages
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := db.Pool.Exec(ctx, query, id)
	return err
}
