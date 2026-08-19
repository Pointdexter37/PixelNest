package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/pixnest/backend/internal/repositories"
	"github.com/pixnest/backend/internal/services"
)

type WallpaperHandler struct {
	service *services.WallpaperService
}

func NewWallpaperHandler(service *services.WallpaperService) *WallpaperHandler {
	return &WallpaperHandler{service: service}
}

func (h *WallpaperHandler) List(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 24)
	if page < 1 || limit < 1 || limit > 100 {
		writeError(w, http.StatusBadRequest, "page must be positive and limit must be between 1 and 100")
		return
	}

	wallpapers, total, err := h.service.List(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list wallpapers")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": wallpapers,
		"meta": map[string]int{"page": page, "limit": limit, "total": total},
	})
}

func (h *WallpaperHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	wallpaper, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, repositories.ErrWallpaperNotFound) {
		writeError(w, http.StatusNotFound, "wallpaper not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get wallpaper")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": wallpaper})
}

func (h *WallpaperHandler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	if err := h.service.RecordDownload(r.Context(), id); errors.Is(err, repositories.ErrWallpaperNotFound) {
		writeError(w, http.StatusNotFound, "wallpaper not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record download")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || value == 0 {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"message": message}})
}
