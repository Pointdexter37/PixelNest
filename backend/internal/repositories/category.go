package repositories

import (
	"context"
	"time"

	"github.com/pixnest/backend/internal/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
}

type MemoryCategoryRepository struct {
	categories []models.Category
}

func NewMemoryCategoryRepository() *MemoryCategoryRepository {
	return &MemoryCategoryRepository{
		categories: []models.Category{{
			ID:        1,
			Name:      "Nature",
			Slug:      "nature",
			CreatedAt: time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
		}},
	}
}

func (r *MemoryCategoryRepository) List(_ context.Context) ([]models.Category, error) {
	return append([]models.Category(nil), r.categories...), nil
}
