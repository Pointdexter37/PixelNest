package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/pixnest/backend/internal/handlers"
	"github.com/pixnest/backend/internal/repositories"
	"github.com/pixnest/backend/internal/services"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	wallpaperRepository := repositories.NewMemoryWallpaperRepository()
	wallpaperService := services.NewWallpaperService(wallpaperRepository)
	wallpaperHandler := handlers.NewWallpaperHandler(wallpaperService)
	categoryHandler := handlers.NewCategoryHandler(repositories.NewMemoryCategoryRepository())

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
	})
	mux.HandleFunc("GET /api/v1/wallpapers", wallpaperHandler.List)
	mux.HandleFunc("GET /api/v1/categories", categoryHandler.List)
	mux.HandleFunc("GET /api/v1/wallpapers/{id}", wallpaperHandler.GetByID)
	mux.HandleFunc("GET /api/v1/wallpapers/{id}/download", wallpaperHandler.Download)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("PixNest API listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}
