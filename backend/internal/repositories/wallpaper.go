package repositories

import (
	"context"
	"errors"

	"github.com/pixnest/backend/internal/models"
)

var ErrWallpaperNotFound = errors.New("wallpaper not found")

type WallpaperRepository interface {
	List(ctx context.Context, offset, limit int, categoryID *int64) ([]models.Wallpaper, error)
	Count(ctx context.Context) (int, error)
	Search(ctx context.Context, query string, offset, limit int, categoryID *int64) ([]models.Wallpaper, int, error)
	GetByID(ctx context.Context, id int64) (models.Wallpaper, error)
	IncrementViews(ctx context.Context, id int64) error
	IncrementDownloads(ctx context.Context, id int64) error
}
