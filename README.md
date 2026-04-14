# Nigerian Institutions API

A comprehensive API and scraper for Nigerian educational institutions, including Universities, Polytechnics, and Colleges of Education. Data is scraped from official sources such as the National Universities Commission (NUC) and the Federal Ministry of Education.

## Project Structure

This project follows the standard Go project layout:

- `cmd/`: Entry points for the application.
  - `api/`: The web API server.
  - `scraper/`: The data scraping utility.
- `internal/`: Private application and library code.
  - `config/`: Configuration management.
  - `constants/`: System-wide constants (e.g., institution types, URLs).
  - `database/`: Database connection and initialization.
  - `dto/`: Data Transfer Objects for API requests and responses.
  - `handlers/`: HTTP request handlers.
  - `model/`: GORM models.
  - `repositories/`: Database access layer.
  - `routes/`: API route definitions.
  - `scraper/`: Scraping logic.
  - `service/`: Business logic layer.
  - `utils/`: Common utilities (e.g., standardized API responses).
- `migrations/`: SQL migration files for database schema management.

## Prerequisites

- [Go](https://go.dev/doc/install) (v1.25.0 or later)
- [PostgreSQL](https://www.postgresql.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for running migrations)
- [Air](https://github.com/cosmtrek/air) (optional, for live reloading)

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/coolpythoncodes/nigerian-universities.git
   cd nigerian-universities
   ```

2. **Configure Environment Variables:**
   Create a `.env` file in the root directory (values are examples):
   ```env
   PORT=8080
   AppEnv=development                 # production disables AutoMigrate
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/nigerian_universities_dev?sslmode=disable

   GOOGLE_CLIENT_ID=...
   GOOGLE_CLIENT_SECRET=...
   GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

   GITHUB_CLIENT_ID=...
   GITHUB_CLIENT_SECRET=...
   GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback

   JWT_SECRET=change-me
   JWT_EXPIRES_IN_HOURS=24
   FRONTEND_URL=http://localhost:3000
   ```

3. **Run Database Migrations (prod / staging):**
   ```bash
   make migrate-up
   ```

4. **(Dev only) AutoMigrate:**
   When `AppEnv != "production"`, the API runs `db.AutoMigrate` for `Institution` and `User` on startup. Use a non-prod database for this (e.g., local Postgres).

5. **Scrape Initial Data:**
   Before running the API, you need to populate the database with institution data:
   ```bash
   go run cmd/scraper/main.go
   ```

6. **Run the API Server:**
   ```bash
   # Using Make (with Air for hot-reload)
   make dev

   # Or directly
   go run cmd/api/main.go
   ```

## Migrations & Atlas
- SQL migrations live in `migrations/`.
- Create a new migration: `make migrate-create name=add_users_table` (fills `migrations/*`).
- Atlas config: `atlas.hcl` points to `DATABASE_URL` (target) and `ATLAS_DEV_URL` (scratch for diffs).
- GitHub Actions workflow `.github/workflows/migrate-and-deploy.yml` runs `atlas migrate apply --env local` on pushes to `main` and can trigger Vercel via `VERCEL_DEPLOY_HOOK_URL`.

## Authentication
- Google OAuth (direct):  
  - `GET /api/v1/auth/google` → redirects to Google  
  - `GET /api/v1/auth/google/callback?code=...` → returns `{ access_token, user }`
- Google (Auth.js flow):  
  - Frontend handles OAuth; call `POST /api/v1/auth/google/login` with `{ id, email, name, avatar_url }` to upsert user and get API JWT.
- GitHub routes are wired (`/api/v1/auth/github`, `/api/v1/auth/github/callback`); implement callback or use an Auth.js-style POST endpoint mirroring the Google flow for GitHub.


## API Endpoints

All API endpoints are prefixed with `/api/v1`.

### Health Check
`GET /api/v1/health`
- Returns the status of the API.

### Institutions
`GET /api/v1/institutions`
- Fetches a list of institutions with support for pagination, searching, and filtering.

**Query Parameters:**
- `page` (default: 1): The page number.
- `limit` (default: 10, max: 100): Number of items per page.
- `search`: Search institutions by name.
- `type`: Filter by institution type. Valid types include:
  - `federal-university`, `state-university`, `private-university`
  - `federal-polytechnic`, `state-polytechnic`, `private-polytechnic`
  - `federal-college-education`, `state-college-education`, `private-college-education`

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/institutions?type=federal-university&search=lagos"
```

**Example Response:**
```json
{
  "success": true,
  "message": "fetched all institutions",
  "data": [
    {
      "id": "...",
      "name": "University of Lagos",
      "vice_chancellor": "...",
      "year_of_establishment": "1962",
      "type": "federal-university",
      "url": "https://unilag.edu.ng"
    }
  ],
  "meta": {
    "current_page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
