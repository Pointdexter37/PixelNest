package services

import (
	"context"

	"github.com/pixnest/backend/internal/models"
	"github.com/pixnest/backend/internal/repositories"
)

type WallpaperService struct {
	repository repositories.WallpaperRepository
}

func NewWallpaperService(repository repositories.WallpaperRepository) *WallpaperService {
	return &WallpaperService{repository: repository}
}

func (s *WallpaperService) List(ctx context.Context, page, limit int) ([]models.Wallpaper, int, error) {
	offset := (page - 1) * limit
	wallpapers, err := s.repository.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return wallpapers, total, nil
}

func (s *WallpaperService) GetByID(ctx context.Context, id int64) (models.Wallpaper, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *WallpaperService) RecordView(ctx context.Context, id int64) error {
	return s.repository.IncrementViews(ctx, id)
}

func (s *WallpaperService) RecordDownload(ctx context.Context, id int64) error {
	return s.repository.IncrementDownloads(ctx, id)
}
