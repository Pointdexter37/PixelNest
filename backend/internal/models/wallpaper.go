package models

import "time"

type Wallpaper struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Slug         string    `json:"slug"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	FileSize     int64     `json:"file_size"`
	Downloads    int64     `json:"downloads"`
	Views        int64     `json:"views"`
	CategoryID   *int64    `json:"category_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
