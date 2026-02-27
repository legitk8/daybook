package services

import (
	"context"
	"daybook/backend/models"
	"daybook/backend/repo"

	"github.com/google/uuid"
)

type PageService struct {
	repo repo.PageRepository
}

func NewPageService(r repo.PageRepository) *PageService {
	return &PageService{
		repo: r,
	}
}

func (s *PageService) CreatePage(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	content string,
) (*models.Page, error) {

	page := &models.Page{
		ID:      uuid.New(),
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	err := s.repo.Create(ctx, page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *PageService) UpdatePage(
	ctx context.Context,
	id uuid.UUID,
	title string,
	content string,
) (*models.Page, error) {

	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	page.Title = title
	page.Content = content

	err = s.repo.Update(ctx, page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *PageService) GetPageByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Page, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PageService) ListPages(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Page, error) {
	return s.repo.GetAllByUser(ctx, userID)
}

func (s *PageService) DeletePage(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repo.SoftDelete(ctx, id)
}
