package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pixelnest/backend/internal/models"
)

type PostgresWallpaperRepository struct {
	db *sql.DB
}

func NewPostgresWallpaperRepository(db *sql.DB) *PostgresWallpaperRepository {
	return &PostgresWallpaperRepository{db: db}
}

func (r *PostgresWallpaperRepository) List(
	ctx context.Context, offset, limit int, categoryID *int64,
) ([]models.Wallpaper, error) {
	return r.queryWallpapers(ctx, "", offset, limit, categoryID)
}

func (r *PostgresWallpaperRepository) Search(
	ctx context.Context, query string, offset, limit int, categoryID *int64,
) ([]models.Wallpaper, int, error) {
	wallpapers, err := r.queryWallpapers(ctx, query, offset, limit, categoryID)
	if err != nil {
		return nil, 0, err
	}
	total, err := r.countMatching(ctx, query, categoryID)
	return wallpapers, total, err
}

func (r *PostgresWallpaperRepository) Count(_ context.Context) (int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM wallpapers").Scan(&total)
	return total, err
}

func (r *PostgresWallpaperRepository) GetByID(ctx context.Context, id int64) (models.Wallpaper, error) {
	row := r.db.QueryRowContext(ctx, wallpaperSelect+" WHERE id = $1", id)
	wallpaper, err := scanWallpaper(row)
	if err == sql.ErrNoRows {
		return models.Wallpaper{}, ErrWallpaperNotFound
	}
	return wallpaper, err
}

func (r *PostgresWallpaperRepository) IncrementViews(ctx context.Context, id int64) error {
	return r.increment(ctx, "views", id)
}

func (r *PostgresWallpaperRepository) IncrementDownloads(ctx context.Context, id int64) error {
	return r.increment(ctx, "downloads", id)
}

const wallpaperSelect = `SELECT id, title, description, slug, image_url,
	thumbnail_url, width, height, file_size, downloads, views, category_id,
	created_at, updated_at FROM wallpapers`

type wallpaperScanner interface {
	Scan(dest ...any) error
}

func scanWallpaper(scanner wallpaperScanner) (models.Wallpaper, error) {
	var wallpaper models.Wallpaper
	err := scanner.Scan(
		&wallpaper.ID, &wallpaper.Title, &wallpaper.Description, &wallpaper.Slug,
		&wallpaper.ImageURL, &wallpaper.ThumbnailURL, &wallpaper.Width,
		&wallpaper.Height, &wallpaper.FileSize, &wallpaper.Downloads,
		&wallpaper.Views, &wallpaper.CategoryID, &wallpaper.CreatedAt,
		&wallpaper.UpdatedAt,
	)
	return wallpaper, err
}

func (r *PostgresWallpaperRepository) queryWallpapers(
	ctx context.Context, query string, offset, limit int, categoryID *int64,
) ([]models.Wallpaper, error) {
	where, args := wallpaperFilters(query, categoryID)
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, wallpaperSelect+" "+where+
		" ORDER BY created_at DESC LIMIT $"+fmt.Sprint(len(args)-1)+
		" OFFSET $"+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallpapers := make([]models.Wallpaper, 0)
	for rows.Next() {
		wallpaper, err := scanWallpaper(rows)
		if err != nil {
			return nil, err
		}
		wallpapers = append(wallpapers, wallpaper)
	}
	return wallpapers, rows.Err()
}

func wallpaperFilters(query string, categoryID *int64) (string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)
	if strings.TrimSpace(query) != "" {
		args = append(args, "%"+strings.TrimSpace(query)+"%")
		conditions = append(conditions, "(title ILIKE $1 OR description ILIKE $1)")
	}
	if categoryID != nil {
		args = append(args, *categoryID)
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *PostgresWallpaperRepository) countMatching(
	ctx context.Context, query string, categoryID *int64,
) (int, error) {
	where, args := wallpaperFilters(query, categoryID)
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wallpapers "+where, args...).Scan(&total)
	return total, err
}

func (r *PostgresWallpaperRepository) increment(ctx context.Context, column string, id int64) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE wallpapers SET "+column+" = "+column+" + 1, updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err == nil && affected == 0 {
		return ErrWallpaperNotFound
	}
	return err
}
