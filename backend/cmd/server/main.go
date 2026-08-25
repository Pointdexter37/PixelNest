package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pixelnest/backend/internal/database"
	"github.com/pixelnest/backend/internal/handlers"
	"github.com/pixelnest/backend/internal/repositories"
	"github.com/pixelnest/backend/internal/services"
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
	wallpaperRepository, closeDatabase, err := newWallpaperRepository()
	if err != nil {
		log.Fatal(err)
	}
	defer closeDatabase()
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

	log.Printf("PixelNest API listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func newWallpaperRepository() (repositories.WallpaperRepository, func(), error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL is not set; using in-memory wallpaper data")
		return repositories.NewMemoryWallpaperRepository(), func() {}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := database.Open(ctx, databaseURL)
	if err != nil {
		return nil, func() {}, err
	}
	log.Println("using PostgreSQL wallpaper repository")
	return repositories.NewPostgresWallpaperRepository(db), func() { db.Close() }, nil
}
