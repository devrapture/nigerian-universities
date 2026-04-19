<img width="71" height="77" alt="image" src="https://github.com/user-attachments/assets/83804a75-30d7-4114-b647-9fa7d7af36f4" />


# Nigerian Institutions API

The Nigerian Institutions API provides a comprehensive and searchable database of educational institutions in Nigeria, including universities, polytechnics, and colleges of education. It consists of a high-performance Go API and a specialized scraper that aggregates data from official government sources (NUC, Federal Ministry of Education).

## Tech Stack

- **Language:** Go (v1.25.0+)
- **Web Framework:** Gin Gonic
- **Database:** PostgreSQL with GORM
- **Migrations:** golang-migrate & Atlas
- **Authentication:** JWT with Google and GitHub OAuth2 integration
- **Documentation:** Swagger (OpenAPI 2.0) via `swaggo/swag`
- **Development Tooling:** Air (hot-reload), Makefile

## Project Structure

The project follows a modular, layered architecture inspired by the standard Go project layout:

- `cmd/`: Application entry points.
  - `api/`: Main API server.
  - `scraper/`: Standalone scraping utility.
- `internal/`: Private core logic.
  - `config/`: Environment-based configuration management.
  - `handlers/`: HTTP interface layer, handling request validation and response formatting.
  - `service/`: Business logic layer, orchestrating repositories and external services.
  - `repositories/`: Data access layer, interacting with PostgreSQL via GORM.
  - `model/`: Database entity definitions.
  - `schema/`: Explicit request/response models for Swagger documentation and API contracts.
  - `middleware/`: Rate limiting, authentication, and logging logic.
  - `scraper/`: Core scraping logic using `colly`.
  - `utils/`: Reusable helpers for JWT, OAuth, and standardized responses.
- `docs/`: Auto-generated Swagger documentation files.

## Environment

Required environment variables:

```env
PORT=8080
AppEnv=development
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

## API Documentation

The API is fully documented using Swagger.

- **Interactive UI:** Available at `/swagger/index.html` when the server is running.
- **Specification:** The raw specification can be found in `docs/swagger.json` and `docs/swagger.yaml`.

## Development Workflows

### 1. Database Management

- **Local Dev:** The API uses `db.AutoMigrate` when `AppEnv` is not `production`.
- **Production/Staging:** Use formal migrations.
  - Create: `make migrate-create name=description`
  - Apply: `make migrate-up`

### 2. Scraping Data

To refresh or populate the database with institutional data, run the scraper:

```bash
go run cmd/scraper/main.go
```

The scraper uses URLs and constants defined in `internal/constants/institution.go`.

### 3. Run the API Server

```bash
# Using Make (with Air for hot-reload)
make dev

# Or directly
go run cmd/api/main.go
```

### 4. Authentication Model

There are two auth mechanisms in the project.

### 1. Bearer JWT

Used for:

- `/api/v1/api-keys/*`

Header:

```http
Authorization: Bearer <jwt>
```

JWTs are returned by the auth endpoints after successful login.

### 2. Product API Key

Used for:

- `/api/v1/institutions`

Header:

```http
X-API-Key: <generated-api-key>
```

The API key is validated against the `product_keys` table through `ProductKeyMiddleware`.

## Rate Limiting

Current route-level throttling in `internal/routes/routes.go`:

- global IP limiter on `/api/v1/*`: `5 req/sec`, burst `10`
- bearer-auth key endpoints: `2 req/sec`, burst `5`

## Endpoints

### Health

`GET /api/v1/health`

Success response:

```json
{
  "status": "ok"
}
```

Failure response:

```json
{
  "status": "db-unreachable",
  "error": "failed to ping database"
}
```

### Auth

#### Google OAuth redirect

`GET /api/v1/auth/google`

Redirects the user to Google OAuth.

#### Google OAuth callback

`GET /api/v1/auth/google/callback?code=...&state=...`

Returns a JWT and user payload on success.

#### Google token login

`POST /api/v1/auth/google/login`

Request body:

```json
{
  "id_token": "google-id-token"
}
```

#### GitHub OAuth redirect

`GET /api/v1/auth/github`

Redirects the user to GitHub OAuth.

#### GitHub OAuth callback

`GET /api/v1/auth/github/callback?code=...&state=...`

Returns a JWT and user payload on success.

#### GitHub token login

`POST /api/v1/auth/github/login`

Request body:

```json
{
  "access_token": "github-access-token"
}
```

Typical auth success shape:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "access_token": "jwt-token",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "john.doe@example.com",
      "name": "John Doe",
      "avatar_url": "https://example.com/avatar.png"
    }
  }
}
```

### Institutions

`GET /api/v1/institutions`

Required header:

```http
X-API-Key: <product-api-key>
```

Query params:

- `page`
- `limit`
- `search`
- `type`

Allowed `type` values:

- `federal-university`
- `state-university`
- `private-university`
- `federal-polytechnic`
- `state-polytechnic`
- `private-polytechnic`
- `federal-college-education`
- `state-college-of-education`
- `private-college-of-education`

Example:

```bash
curl -H "X-API-Key: sk_live_..." \
  "http://localhost:8080/api/v1/institutions?type=federal-university&search=lagos&page=1&limit=10"
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

### API Keys

These endpoints require Bearer auth.

#### Generate key

`POST /api/v1/api-keys/generate`

Header:

```http
Authorization: Bearer <jwt>
```

#### List keys

`GET /api/v1/api-keys?page=1&limit=10`

#### Revoke key

`POST /api/v1/api-keys/{key_id}/revoke`

## Swagger

Swagger UI is served from:

```text
http://localhost:8080/swagger/index.html
```

### Regenerate docs

Use this command:

```bash
swag init -g main.go -d cmd/api,internal/handlers,internal/routes,internal/schema,internal/constants --parseInternal
```

Why these directories are included:

- `cmd/api`: general Swagger metadata
- `internal/handlers`: auth, institutions, keys annotations
- `internal/routes`: health endpoint annotations
- `internal/schema`: Swagger-only request/response models
- `internal/constants`: institution type resolution

### Authorization in Swagger

Swagger is configured with a global Bearer auth scheme for JWT-protected routes.

Use the **Authorize** button for:

- `/api/v1/api-keys/`*

For institutions, Swagger still expects manual `X-API-Key` input because that endpoint uses product keys, not bearer tokens.

## Notes

- `institutions` is protected by product API keys, not JWT.
- `api-keys` endpoints are protected by JWT, not `X-API-Key`.
- health is public.
- auth endpoints are public.
- Swagger schemas for docs are intentionally separated into `internal/schema`.

