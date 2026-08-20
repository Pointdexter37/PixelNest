package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/pixnest/backend/internal/models"
)

type MemoryWallpaperRepository struct {
	mu         sync.RWMutex
	wallpapers []models.Wallpaper
}

func NewMemoryWallpaperRepository() *MemoryWallpaperRepository {
	// Seed one record so the API and frontend can be developed before PostgreSQL is connected.
	return &MemoryWallpaperRepository{
		wallpapers: []models.Wallpaper{
			{
				ID:           1,
				Title:        "Mountain Lake",
				Description:  "A calm mountain lake beneath a clear evening sky.",
				Slug:         "mountain-lake",
				ImageURL:     "https://images.unsplash.com/photo-1500534623283-312aade485b7",
				ThumbnailURL: "https://images.unsplash.com/photo-1500534623283-312aade485b7?w=800",
				Width:        1920,
				Height:       1080,
				FileSize:     2457600,
				CategoryID:   nil,
				CreatedAt:    time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func (r *MemoryWallpaperRepository) List(_ context.Context, offset, limit int) ([]models.Wallpaper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if offset >= len(r.wallpapers) {
		return []models.Wallpaper{}, nil
	}
	end := offset + limit
	if end > len(r.wallpapers) {
		end = len(r.wallpapers)
	}
	return append([]models.Wallpaper(nil), r.wallpapers[offset:end]...), nil
}

func (r *MemoryWallpaperRepository) Count(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.wallpapers), nil
}

func (r *MemoryWallpaperRepository) GetByID(_ context.Context, id int64) (models.Wallpaper, error) {
	r.mu.RLock() // multiple goroutines can safely read at the same time.
	defer r.mu.RUnlock()

	for _, wallpaper := range r.wallpapers {
		if wallpaper.ID == id {
			return wallpaper, nil
		}
	}
	return models.Wallpaper{}, ErrWallpaperNotFound
}

func (r *MemoryWallpaperRepository) IncrementDownloads(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.wallpapers {
		if r.wallpapers[index].ID == id {
			r.wallpapers[index].Downloads++
			return nil
		}
	}
	return ErrWallpaperNotFound
}
