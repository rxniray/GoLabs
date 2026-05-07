# Project Specification

## Overview

A monolithic web application written in Go, exposing a REST API built with the Fiber framework. The service provides user authentication, image uploads, PDF document generation, and transactional email sending. The emphasis is on simplicity: one binary, one database, minimal dependencies.

## Tech Stack

| Concern          | Choice                                    |
| ---------------- | ----------------------------------------- |
| Language         | Go (1.22+)                                |
| HTTP framework   | [Fiber v2](https://gofiber.io)            |
| Database         | PostgreSQL 15+                            |
| DB driver        | `pgx` (via `database/sql` or native)      |
| Migrations       | [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) |
| Authentication   | JWT (access tokens) + bcrypt for passwords |
| Image storage    | Local filesystem                          |
| PDF generation   | [`gofpdf`](https://github.com/jung-kurt/gofpdf) or [`maroto`](https://github.com/johnfercher/maroto) |
| Email            | SMTP via `net/smtp` or `go-mail/mail`     |
| Configuration    | Environment variables (via `godotenv` for local dev) |
| Logging          | `log/slog` (standard library)             |

## Architecture

Monolithic. Single binary, single deployment unit. Code is organized by feature (module) rather than by technical layer.

```
/cmd
  /api                  # main.go — entrypoint, wiring, server bootstrap
/internal
  /config               # env loading, config struct
  /database             # pgx pool setup, health check
  /migrations           # .sql migration files (embedded)
  /auth                 # sign-up, sign-in, JWT middleware
  /users                # user model, repository
  /images               # upload handler, storage on disk
  /documents            # PDF generation handler
  /email                # SMTP client wrapper, templates
  /server               # router setup, middleware, error handling
/pkg                    # (optional) reusable helpers, if any
/storage
  /uploads              # image files (gitignored)
/migrations             # *.up.sql / *.down.sql
.env.example
Dockerfile
docker-compose.yml
go.mod
```

Each module under `/internal` exposes a small surface: a `Handler` (HTTP), a `Service` (business logic), and a `Repository` (DB access) where applicable. Modules talk to each other through service interfaces, not through HTTP.

## Modules

### 1. Authentication (`/internal/auth`)

**Endpoints**

| Method | Path                | Description                           | Auth |
| ------ | ------------------- | ------------------------------------- | ---- |
| POST   | `/api/auth/signup`  | Register a new user                   | No   |
| POST   | `/api/auth/signin`  | Authenticate, returns JWT             | No   |
| GET    | `/api/auth/me`      | Returns current user profile          | Yes  |

**Sign-up flow**
1. Accept `email`, `password`, `name`.
2. Validate (valid email, password ≥ 8 chars).
3. Hash password with bcrypt (cost 10).
4. Insert user; return `201 Created`.
5. Send welcome email via the email module (async, fire-and-forget).

**Sign-in flow**
1. Look up user by email.
2. Compare bcrypt hash.
3. Issue a JWT signed with HMAC-SHA256, claims: `sub` (user id), `exp` (24h), `iat`.
4. Return `{ "token": "..." }`.

**JWT middleware**
Extracts `Authorization: Bearer <token>`, validates signature and expiry, puts `userID` in Fiber's `c.Locals("userID", ...)`.

### 2. Users (`/internal/users`)

Owns the `users` table. Exposes a repository used by the auth module. No public endpoints beyond those in auth.

### 3. Images (`/internal/images`)

**Endpoints**

| Method | Path                      | Description                      | Auth |
| ------ | ------------------------- | -------------------------------- | ---- |
| POST   | `/api/images`             | Upload an image (multipart form) | Yes  |
| GET    | `/api/images/:id`         | Download/serve image             | Yes  |
| DELETE | `/api/images/:id`         | Delete image                     | Yes  |

**Upload rules**
- Field name: `file`.
- Allowed MIME types: `image/jpeg`, `image/png`, `image/webp`.
- Max size: 10 MB (configurable via `MAX_UPLOAD_SIZE_MB`).
- Filename is regenerated as `<uuid>.<ext>` to avoid collisions and path traversal.
- Stored under `./storage/uploads/<year>/<month>/<uuid>.<ext>`.
- Metadata (id, owner_id, original_name, stored_path, mime, size, created_at) saved in `images` table.

### 4. Documents (`/internal/documents`)

**Endpoints**

| Method | Path                       | Description                         | Auth |
| ------ | -------------------------- | ----------------------------------- | ---- |
| POST   | `/api/documents/generate`  | Generate a PDF from a template name + JSON data; returns the PDF binary | Yes  |

**Request body**
```json
{
  "template": "invoice",
  "data": { "number": "INV-001", "amount": 100.00, "client": "Acme" }
}
```

**Behavior**
- Template is looked up by name from a registry (Go code). Start with 1–2 templates (e.g., `invoice`, `report`).
- PDF is rendered in-memory and returned with `Content-Type: application/pdf`.
- No persistence in the first version — callers save the bytes themselves.

### 5. Email (`/internal/email`)

Internal module, no public endpoints. Exposes `Send(to, subject, body string) error` and `SendTemplate(to, templateName string, data any) error`.

**Config (env)**
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM`

**Templates**
Stored as Go `html/template` files under `/internal/email/templates/`. Start with: `welcome.html`, `password_reset.html` (reset flow can be stubbed for v1).

For local development, point SMTP at [Mailtrap](https://mailtrap.io) or [MailHog](https://github.com/mailhog/MailHog) — no real mail is sent.

## Database

### Schema (initial)

**`users`**
| Column         | Type          | Notes                     |
| -------------- | ------------- | ------------------------- |
| id             | UUID          | PK, `gen_random_uuid()`   |
| email          | TEXT          | UNIQUE, NOT NULL          |
| password_hash  | TEXT          | NOT NULL                  |
| name           | TEXT          | NOT NULL                  |
| created_at     | TIMESTAMPTZ   | DEFAULT `now()`           |
| updated_at     | TIMESTAMPTZ   | DEFAULT `now()`           |

**`images`**
| Column         | Type          | Notes                     |
| -------------- | ------------- | ------------------------- |
| id             | UUID          | PK                        |
| owner_id       | UUID          | FK → users(id)            |
| original_name  | TEXT          |                           |
| stored_path    | TEXT          | NOT NULL                  |
| mime_type      | TEXT          |                           |
| size_bytes     | BIGINT        |                           |
| created_at     | TIMESTAMPTZ   | DEFAULT `now()`           |

### Migrations

- Files live in `/migrations` as pairs: `0001_init.up.sql` / `0001_init.down.sql`.
- Embedded into the binary using `//go:embed` and run at startup via `golang-migrate`.
- Running up migrations is automatic on boot; down migrations are only run via CLI (`./api migrate down`).

## Configuration

All config via environment variables. `.env.example` committed, `.env` gitignored.

```
# Server
APP_PORT=8080
APP_ENV=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/appdb?sslmode=disable

# Auth
JWT_SECRET=change-me-in-production
JWT_TTL_HOURS=24

# Uploads
UPLOAD_DIR=./storage/uploads
MAX_UPLOAD_SIZE_MB=10

# SMTP
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=no-reply@example.com
```

## Error Handling

A single Fiber error handler returns consistent JSON:

```json
{ "error": { "code": "INVALID_CREDENTIALS", "message": "Email or password is incorrect" } }
```

Standard HTTP codes: `400` validation, `401` auth, `403` forbidden, `404` not found, `409` conflict (e.g., email already taken), `413` payload too large, `500` server.

## Logging

`slog` with JSON handler in production, text in development. Request-level middleware logs method, path, status, latency, and (if authenticated) user id.

## Running Locally

```bash
# 1. Start Postgres + MailHog
docker-compose up -d

# 2. Copy env
cp .env.example .env

# 3. Run migrations (automatic on boot, or manually:)
go run ./cmd/api migrate up

# 4. Start the server
go run ./cmd/api
```

Server available at `http://localhost:8080`.

## Out of Scope (v1)

Intentionally excluded to keep things simple:
- Refresh tokens / token revocation
- Role-based access control
- Password reset emails (stub only)
- Rate limiting
- Horizontal scaling / session store
- Cloud object storage (can be added behind the images module's storage interface)
- Document persistence / async generation queue
- API versioning beyond `/api`
- OpenAPI/Swagger generation

## Future Extensions

The module boundaries are designed so these additions don't require refactoring:
- Swap local image storage for S3 by implementing a `Storage` interface.
- Add DOCX generation as a second template engine in the documents module.
- Extract the email module into a worker if throughput grows.
