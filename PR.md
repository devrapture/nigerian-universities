# Pull Request Template

## Description

<!-- Clearly and concisely describe the changes you made and the motivation behind them. -->

- Summary:
  - What changed:
    - Upgraded Go version from 1.21 to 1.23.
    - Implemented database configuration and connection logic using GORM and PostgreSQL.
    - Added environment variable management using `godotenv`.
    - Integrated `Air` for live-reloading during development.
    - Created a `makefile` to streamline migrations and development workflows.
    - Scaffolded the `University` model and initial SQL migrations.
  - Why:
    - To provide a persistent storage layer for university data.
    - To improve development productivity and maintainability.
    - To modernize the Go toolchain.
  - How (high-level):
    - Added `internal/config` to load settings from environment variables.
    - Added `internal/database` to manage PostgreSQL connections with pooling and conditional auto-migration for non-production environments.
    - Added `.air.toml` and `makefile` for automation.
- Affected areas/modules:
  - Application entry point (`cmd/api/main.go`).
  - Configuration management.
  - Database layer.
- Risks/Trade-offs:
  - Dependency on external PostgreSQL database.
  - The `ConnectDB` logic currently uses `AutoMigrate` in development, which might lead to schema drift if not managed carefully alongside SQL migrations.
- Rollback plan:
  - Revert the changes to `cmd/api/main.go` and `go.mod`.
  - Use `make migrate-down` if migrations were applied.

---

## References

<!-- Add links to the related ticket, issue, Slack conversation, or any other relevant references. -->

- Ticket/Issue: N/A
- Other references: N/A

---

## Screenshots/Recordings

### Before

<!-- Add screenshots or screen recordings of the behavior before the change. -->

N/A (Backend infrastructure changes)

### After

<!-- Add screenshots or screen recordings of the behavior after the change. -->

N/A (Backend infrastructure changes)

---

## Testing Steps

<!-- Provide a step-by-step guide on how to test your changes. -->

1. Ensure PostgreSQL is running and you have a database created.
2. Create a `.env` file in the root with `DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable`.
3. Run `make migrate-up` to verify the migration tool works (even with empty migrations).
4. Run `make dev` to start the application with live reloading.
5. Check the logs to confirm the database connection was successfully established and the server started on the configured port.

---

## Checklist

<!-- Check off each item to ensure your PR is ready for review. -->

- [x] I have run the application to ensure my changes work as expected.
- [x] I have added tests where applicable.
- [ ] I have included screenshots or recordings of any UI changes.
- [x] I have included all necessary references to tickets, issues, or conversations.
- [x] This PR is ready for review and meets the contribution guidelines.

---

## Notes

<!-- Include any additional comments, concerns, or questions for the reviewer. -->

- The `internal/model/university.go` uses a singular `University` with UUID, while the existing `models/university_model.go` uses a plural `Universities`. The `ConnectDB` currently uses the latter for auto-migration. This should be unified in a follow-up PR.
- Migrations in the `migrations` folder are currently empty placeholders.
