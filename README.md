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
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
   App_Env=development
   ```

3. **Run Database Migrations:**
   ```bash
   make migrate-up
   ```

4. **Scrape Initial Data:**
   Before running the API, you need to populate the database with institution data:
   ```bash
   go run cmd/scraper/main.go
   ```

5. **Run the API Server:**
   ```bash
   # Using Make (with Air for hot-reload)
   make dev

   # Or directly
   go run cmd/api/main.go
   ```

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
