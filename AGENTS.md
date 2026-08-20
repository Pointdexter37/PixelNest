# Wallpaper Platform — Agent Guide

## Project Overview

Build a full-stack website where users can discover, preview, and download high-quality desktop wallpapers for free.

The project should be designed as a real production-style application rather than a basic CRUD project. Prioritize performance, clean architecture, scalability, SEO, and a polished user experience.

## Tech Stack

### Frontend
- TypeScript
- Next.js
- React
- Tailwind CSS
- shadcn/ui

### Backend
- Go
- REST API
- PostgreSQL
- Redis where caching/rate limiting is useful

### Storage & Infrastructure
- Cloudflare R2 for wallpaper/image storage
- CDN-friendly image delivery
- Docker for local development and deployment
- GitHub Actions for CI/CD

## Core Features

### User Features
- Browse wallpapers
- Search wallpapers
- Filter by category
- Filter by resolution/aspect ratio
- View wallpaper details
- Download wallpapers
- Track views and downloads
- Responsive desktop/mobile interface

### Wallpaper Metadata
Each wallpaper should support:
- Title
- Description
- Category
- Tags
- Resolution
- Aspect ratio
- File size
- Image URL
- Thumbnail URL
- View count
- Download count
- Created date

### Admin Features
Create a protected admin area where administrators can:
- Upload wallpapers
- Generate/store thumbnails
- Add and edit metadata
- Create/manage categories
- Manage tags
- Delete wallpapers
- View basic download/view statistics

## Suggested Categories

- Anime
- Nature
- Cars
- Gaming
- Minimal
- Space
- Cyberpunk
- Abstract
- Architecture
- Technology
- Movies
- Animals

## Backend Architecture

Use Go as a separate backend service.

Suggested structure:

backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   ├── storage/
│   └── routes/
├── migrations/
├── Dockerfile
└── go.mod

Keep business logic out of HTTP handlers. Use a layered architecture:

HTTP Handler → Service → Repository → PostgreSQL

Storage operations should be handled separately through a storage interface so Cloudflare R2 can be replaced later if necessary.

## Frontend Architecture

Suggested structure:

frontend/
├── src/
│   ├── app/
│   ├── components/
│   ├── lib/
│   ├── hooks/
│   ├── services/
│   ├── types/
│   └── styles/
├── public/
├── Dockerfile
└── package.json

The frontend communicates with the Go backend through REST APIs.

Do not put database access directly inside the frontend.

## API Design

Use RESTful endpoints.

Example:

GET    /api/v1/wallpapers
GET    /api/v1/wallpapers/:id
GET    /api/v1/wallpapers/:id/download
GET    /api/v1/categories
GET    /api/v1/tags
GET    /api/v1/search

Admin:

POST   /api/v1/admin/wallpapers
PATCH  /api/v1/admin/wallpapers/:id
DELETE /api/v1/admin/wallpapers/:id
POST   /api/v1/admin/categories
PATCH  /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id

Return consistent JSON responses and proper HTTP status codes.

## Database

Use PostgreSQL.

Initial entities:

### wallpapers
- id
- title
- description
- slug
- image_url
- thumbnail_url
- width
- height
- file_size
- downloads
- views
- category_id
- created_at
- updated_at

### categories
- id
- name
- slug
- created_at

### tags
- id
- name
- slug

### wallpaper_tags
- wallpaper_id
- tag_id

Add indexes for fields frequently used in search/filtering.

Use database migrations rather than manually modifying production databases.

## Image Storage

Use Cloudflare R2 for original wallpaper files and thumbnails.

Recommended flow:

1. Admin requests an upload URL from the Go backend.
2. Backend validates the request.
3. Upload goes to R2.
4. Backend stores the resulting metadata in PostgreSQL.
5. Frontend displays optimized thumbnails.
6. Downloads use the original/high-resolution file.

Avoid storing large image files directly in PostgreSQL.

## Image Processing

For uploaded images:

- Validate file type.
- Validate maximum dimensions/file size.
- Detect width and height.
- Generate thumbnails.
- Prefer modern formats such as WebP/AVIF where appropriate.
- Preserve the original high-resolution file for downloads.

Do not trust file extensions alone. Validate the actual image content.

## Search

Start with PostgreSQL-based search.

Support:
- Title search
- Tag search
- Category filtering
- Resolution filtering
- Aspect-ratio filtering

Design the repository layer so a dedicated search engine can be introduced later if the project grows.

## Performance

Performance is a major project goal.

Frontend:
- Use Next.js image optimization where appropriate.
- Lazy-load images.
- Use thumbnails for gallery pages.
- Use pagination or infinite scrolling.
- Avoid loading full-resolution images in galleries.

Backend:
- Add database indexes.
- Cache popular queries with Redis when useful.
- Use connection pooling.
- Avoid N+1 database queries.

Storage:
- Serve images through a CDN.
- Cache static assets aggressively.

## Security

Implement:
- Admin authentication
- Authorization middleware
- Input validation
- Upload restrictions
- File size limits
- MIME/content validation
- Rate limiting for public APIs
- Rate limiting for downloads where appropriate
- CORS configuration
- Secure environment variables

Never commit secrets, API keys, database credentials, or R2 credentials to Git.

## SEO

SEO is important because wallpaper sites depend heavily on search traffic.

Implement:
- SEO-friendly slugs
- Dynamic metadata
- Open Graph metadata
- Sitemap
- Robots.txt
- Canonical URLs
- Descriptive image alt text
- Structured metadata where useful

Example URL:

/wallpaper/cyberpunk-city-4k

Avoid exposing meaningless URLs such as:

/wallpaper/183728

## UI Requirements

The design should be modern, minimal, and image-focused.

Homepage should contain:
- Header/navigation
- Search bar
- Featured/trending wallpapers
- Categories
- Latest wallpapers
- Responsive wallpaper grid

Wallpaper detail page:
- Large preview
- Title
- Resolution
- File size
- Tags
- Category
- Download button
- Related wallpapers

Use skeleton loading states and useful empty/error states.

## Development Rules

- Use TypeScript strictly on the frontend.
- Use idiomatic Go on the backend.
- Keep frontend and backend independently deployable.
- Keep secrets in environment variables.
- Use meaningful commit messages.
- Write reusable components.
- Avoid unnecessary dependencies.
- Prefer simple solutions before introducing complex infrastructure.
- Validate data at API boundaries.
- Return useful errors without exposing internal implementation details.

## Environment Variables

Frontend example:

NEXT_PUBLIC_API_URL=

Backend example:

PORT=
DATABASE_URL=
REDIS_URL=
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET_NAME=
R2_PUBLIC_URL=
ADMIN_JWT_SECRET=

Never commit `.env` files containing real credentials.

## Local Development

Recommended structure:

wallpaper-platform/
├── frontend/
├── backend/
├── docker-compose.yml
├── .gitignore
└── README.md

Docker Compose can run:
- PostgreSQL
- Redis
- Go backend
- Next.js frontend

## Testing

Backend:
- Unit tests for services
- Repository tests where useful
- HTTP handler tests
- Upload validation tests

Frontend:
- Component tests for important UI
- End-to-end tests for browsing and downloading

Important flows to test:
1. Browse wallpapers.
2. Search.
3. Filter.
4. Open wallpaper details.
5. Download wallpaper.
6. Admin authentication.
7. Admin upload.
8. Delete wallpaper.

## Deployment

Target architecture:

User
  ↓
Next.js Frontend
  ↓
Go REST API
  ↓
PostgreSQL

Go API
  ↓
Cloudflare R2

Cloudflare/CDN
  ↓
Wallpaper assets

Use Docker for backend deployment and GitHub Actions for CI/CD.

## Development Phases

### Phase 1 — Foundation
- Initialize frontend
- Initialize Go backend
- Setup PostgreSQL
- Create migrations
- Establish frontend/backend communication

### Phase 2 — Core Wallpaper System
- Wallpaper database model
- Categories
- Tags
- Wallpaper listing
- Wallpaper detail page
- Download endpoint

### Phase 3 — Admin
- Admin authentication
- Upload workflow
- R2 integration
- Image validation
- Thumbnail generation
- Wallpaper management

### Phase 4 — Search & Discovery
- Search
- Filtering
- Sorting
- Related wallpapers
- Trending wallpapers

### Phase 5 — Production Hardening
- Redis caching
- Rate limiting
- SEO
- Error handling
- Logging
- Testing
- Image optimization
- CDN configuration

### Phase 6 — Deployment
- Dockerize services
- Configure production database
- Configure R2
- Setup CI/CD
- Deploy frontend
- Deploy backend
- Configure domain and HTTPS

## Future Features

Do not implement these before the core platform is stable:

- User accounts
- Favorites
- User uploads
- Collections
- AI-powered tagging
- Personalized recommendations
- Color-based wallpaper search
- Wallpaper packs/ZIP downloads
- Analytics dashboard
- AI wallpaper generation

## Definition of Done

The MVP is complete when:

- Users can browse wallpapers.
- Users can search and filter wallpapers.
- Users can open wallpaper detail pages.
- Users can download wallpapers.
- Admins can securely upload wallpapers.
- Images are stored in Cloudflare R2.
- Metadata is stored in PostgreSQL.
- Thumbnails are used for gallery views.
- The frontend and Go backend communicate through REST APIs.
- Basic SEO is implemented.
- The application works responsively.
- The project can be deployed using Docker.
- No secrets are committed to Git.

## Engineering Priority

When making implementation decisions, prioritize in this order:

1. Correctness
2. Security
3. Performance
4. Maintainability
5. Scalability
6. Developer experience

Do not over-engineer the MVP. Build the simplest production-quality version first, then optimize based on actual needs.

## Project Workflow

Progress should be deliberate and incremental. We are moving slowly, one step at a time, and completing and reviewing each focused change before starting the next one. After every step, the user will review the implementation and read the code before additional work continues.

### Incremental Change Limits

- Move to the next focused step only after the current step is complete and reviewable.
- Do not generate or modify more than 100 lines of code in a single step.
- If a change requires more than 100 lines, split it into smaller sequential steps and pause for review between them.
- Keep each step narrowly scoped; avoid bundling unrelated refactors or features.

### Learning Focus

This project is being built as a learning exercise as well as a working wallpaper platform.

- Prefer clear, conventional implementations over clever or overly abstract solutions.
- Explain important architectural decisions and non-obvious code with concise comments or documentation.
- Keep comments focused on why the code exists, the tradeoffs involved, or how the pieces connect.
- Introduce new tools and patterns gradually so each step can be understood and reviewed before moving on.
