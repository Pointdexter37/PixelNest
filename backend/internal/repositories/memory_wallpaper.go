package repositories

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pixelnest/backend/internal/models"
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
				CategoryID:   int64Pointer(1),
				CreatedAt:    time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}

func (r *MemoryWallpaperRepository) List(_ context.Context, offset, limit int, categoryID *int64) ([]models.Wallpaper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matches := filterByCategory(r.wallpapers, categoryID)
	if offset >= len(matches) {
		return []models.Wallpaper{}, nil
	}
	end := offset + limit
	if end > len(matches) {
		end = len(matches)
	}
	return append([]models.Wallpaper(nil), matches[offset:end]...), nil
}

func (r *MemoryWallpaperRepository) Count(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.wallpapers), nil
}

func (r *MemoryWallpaperRepository) Search(_ context.Context, query string, offset, limit int, categoryID *int64) ([]models.Wallpaper, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	matches := make([]models.Wallpaper, 0)
	for _, wallpaper := range filterByCategory(r.wallpapers, categoryID) {
		if strings.Contains(strings.ToLower(wallpaper.Title), query) ||
			strings.Contains(strings.ToLower(wallpaper.Description), query) {
			matches = append(matches, wallpaper)
		}
	}
	if offset >= len(matches) {
		return []models.Wallpaper{}, len(matches), nil
	}
	end := offset + limit
	if end > len(matches) {
		end = len(matches)
	}
	return append([]models.Wallpaper(nil), matches[offset:end]...), len(matches), nil
}

func filterByCategory(wallpapers []models.Wallpaper, categoryID *int64) []models.Wallpaper {
	if categoryID == nil {
		return wallpapers
	}
	matches := make([]models.Wallpaper, 0)
	for _, wallpaper := range wallpapers {
		if wallpaper.CategoryID != nil && *wallpaper.CategoryID == *categoryID {
			matches = append(matches, wallpaper)
		}
	}
	return matches
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

func (r *MemoryWallpaperRepository) IncrementViews(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.wallpapers {
		if r.wallpapers[index].ID == id {
			r.wallpapers[index].Views++
			return nil
		}
	}
	return ErrWallpaperNotFound
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
