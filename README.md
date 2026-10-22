# API URL Shortener

A REST API for creating and managing short links. Users can register, log in, create links with either a custom or generated slug, browse and search their links, and delete them. Visiting an active short link redirects to its original URL.

## Features

- User registration, login, and logout.
- JWT-based authentication.
- Short links with custom or randomly generated slugs.
- Search and pagination for a user's links.
- Soft deletion of links.
- Redirects from short slugs to original URLs.

## Technology Stack

- **Go** — version specified in `go.mod`
- **Gin** — HTTP framework
- **PostgreSQL** — stores users and links
- **Redis** — stores revoked tokens after logout
- **JWT** — authentication tokens
- **golang-migrate** — database migrations

## Prerequisites

- Go version specified in `go.mod`
- PostgreSQL
- Redis
- `golang-migrate` CLI to apply database migrations
- GNU Make to use the commands in the `makefile`

## Run Locally

1. Clone the repository and change into its directory:

   ```bash
   git clone https://github.com/rcontrisha/api-url-shortener.git
   cd api-url-shortener
   ```

2. Create a `.env` file from the example:

   **PowerShell**
   ```powershell
   Copy-Item .env.example .env
   ```

   **Linux/macOS**
   ```bash
   cp .env.example .env
   ```

3. Configure the values in `.env`:

   | Variable | Description |
   | --- | --- |
   | `HOST` | Application bind address, e.g. `0.0.0.0` |
   | `PORT` | HTTP port, e.g. `8080` |
   | `DB_USER` | PostgreSQL username |
   | `DB_PASSWORD` | PostgreSQL password |
   | `DB_HOST` | PostgreSQL host |
   | `DB_PORT` | PostgreSQL port |
   | `DB_NAME` | PostgreSQL database name |
   | `REDIS_USER` | Redis username, if configured |
   | `REDIS_PASSWORD` | Redis password, if configured |
   | `REDIS_HOST` | Redis host |
   | `REDIS_PORT` | Redis port |
   | `JWT_KEY` | Secret used to sign JWTs |
   | `JWT_ISSUER` | JWT issuer claim |

   Replace the placeholders in `.env.example` with values for your local services. Keep `.env` and `JWT_KEY` private; do not commit secrets.

4. Apply the database migrations:

   ```bash
   make migrate-up
   ```

   This requires GNU Make, `golang-migrate`, and a valid database configuration in `.env`. The application does not apply migrations automatically at startup.

5. Download dependencies and start the API:

   ```bash
   go mod download
   go run .
   ```

The server uses `HOST` and `PORT` from `.env`. For example, if configured with `HOST=localhost` and `PORT=8080`, the local base URL is `http://localhost:8080`.

## API Documentation

All endpoints return JSON except the short-link redirect endpoint. The standard JSON response shape is:

```json
{
  "success": true,
  "message": "Operation completed.",
  "results": null
}
```

### Register

`POST /auth/register`

Request:

```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

Example:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"your-password"}'
```

Success response: `201 Created`.

### Log in

`POST /auth/login`

Use the same `email` and `password` request fields as registration.

Example:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"your-password"}'
```

Successful response (`200 OK`) includes a token:

```json
{
  "success": true,
  "message": "Login Success.",
  "results": {
    "token": "<jwt>"
  }
}
```

The JWT expires after 60 minutes.

### Log out

`POST /auth/logout`

Requires the `Authorization: Bearer <jwt>` header.

Example:

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer <jwt>"
```

Success response: `200 OK`. The token is added to the Redis blacklist until it expires.

### Create a short link

`POST /link/create`

Requires the `Authorization: Bearer <jwt>` header.

Request with an automatically generated slug:

```json
{
  "original_url": "https://example.com/articles/getting-started"
}
```

Request with a custom slug:

```json
{
  "original_url": "https://example.com/articles/getting-started",
  "slug": "getting-started"
}
```

Example:

```bash
curl -X POST http://localhost:8080/link/create \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://example.com/articles/getting-started","slug":"getting-started"}'
```

`original_url` is required and must be a valid URL. Generated slugs contain 6 alphanumeric characters. Custom slugs must be 3–50 characters and cannot use reserved words: `api`, `login`, `register`, or `dashboard`.

Success response: `201 Created`.

### List a user's links

`GET /link/links`

Requires the `Authorization: Bearer <jwt>` header.

| Query parameter | Default | Description |
| --- | --- | --- |
| `search` | Empty | Searches slugs and original URLs |
| `page` | `1` | Results page |
| `limit` | `10` | Results per page |

Example:

```bash
curl "http://localhost:8080/link/links?search=example&page=1&limit=10" \
  -H "Authorization: Bearer <jwt>"
```

Example response:

```json
{
  "success": true,
  "message": "successfully retrieve user links",
  "results": {
    "data": [
      {
        "id": 1,
        "slug": "getting-started",
        "original_url": "https://example.com/articles/getting-started",
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "meta": {
      "total_active": 1,
      "page": 1,
      "limit": 10,
      "total_pages": 1
    }
  }
}
```

### Delete a short link

`DELETE /link/:id/delete`

Requires the `Authorization: Bearer <jwt>` header. Users can only delete their own links.

Example:

```bash
curl -X DELETE http://localhost:8080/link/1/delete \
  -H "Authorization: Bearer <jwt>"
```

Deletion is a soft delete: the link remains in the database but is no longer active. Success response: `200 OK`.

### Redirect to the original URL

`GET /:slug`

This endpoint does not require authentication. If the slug is active, the API returns a `301` redirect to the original URL.

Example:

```bash
curl -i http://localhost:8080/getting-started
```

## Assumptions and Design Decisions

- PostgreSQL stores user and link data; Redis stores revoked JWTs.
- Passwords are stored as Argon2id hashes, not plaintext.
- JWTs use HS256 and expire after 60 minutes.
- Deleted links use soft deletion.
- Slugs must be unique among active links.
- The `total_active` pagination value counts all active links belonging to the user; it is not filtered by the `search` query.
- The custom-slug validator currently accepts a specific character pattern. Its error message describes the allowed characters more broadly; see `internal/utils/slug-validator.go` for the implementation.

## Database Migrations

Migration files are in `db/migrations`. The `makefile` provides:

```bash
make migrate-up
make migrate-down
```

`migrate-down` rolls back the most recent migration; use it with care.
