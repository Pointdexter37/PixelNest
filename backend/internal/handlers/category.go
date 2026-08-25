package handlers

import (
	"net/http"

	"github.com/pixnest/backend/internal/repositories"
)

type CategoryHandler struct {
	repository repositories.CategoryRepository
}

func NewCategoryHandler(repository repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repository: repository}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repository.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": categories})
}
