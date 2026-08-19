# PixNest

PixNest is a wallpaper platform with a Next.js frontend and Go REST API.

## Local development

1. Copy `.env.example` to `.env`.
2. Start local infrastructure with `docker compose up -d`.
3. Run the API with `cd backend && go run ./cmd/server`.
4. Install frontend dependencies with `cd frontend && npm install`.
5. Start the frontend with `npm run dev`.

The API health check is available at `http://localhost:8080/healthz`.

Development proceeds incrementally. Each focused step should be reviewed before the next implementation step begins.

