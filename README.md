# Linkboard

**A backend-first bookmark management REST API built with Go and PostgreSQL.**

Linkboard is a learning-focused backend engineering project designed to explore real-world backend concepts, database design, API development, authentication, security, and system architecture.

The project follows a layered architecture using Go's standard `net/http` package and PostgreSQL through `pgx`.

## Tech Stack

| Technology | Purpose |
|---|---|
| Go | Backend development |
| `net/http` | HTTP server and routing |
| PostgreSQL 17 | Relational database |
| pgx / pgxpool | PostgreSQL driver and connection pooling |
| SQL | Database queries and migrations |
| Docker Compose | Local database environment |
| Git & GitHub | Version control |

## Architecture

Linkboard follows a layered backend architecture:

```text
Client
  |
  v
HTTP Router
  |
  v
Handler
  |
  v
Service
  |
  v
Repository
  |
  v
PostgreSQL
```

### Layer responsibilities

- **Handler:** Handles HTTP requests, parses input and sends responses.
- **Service:** Contains business logic and application rules.
- **Repository:** Handles database operations and SQL queries.
- **Database:** Stores application data in PostgreSQL.

## Current Features

### Bookmark Management

- Create bookmarks.
- Retrieve all bookmarks.
- Retrieve a bookmark by ID.
- Update existing bookmarks.
- Soft-delete bookmarks.
- Pagination using limit and offset.
- PostgreSQL-backed persistence.

### API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/bookmarks` | List bookmarks |
| GET | `/bookmarks/{id}` | Get a bookmark |
| POST | `/bookmarks` | Create a bookmark |
| PUT | `/bookmarks/{id}` | Update a bookmark |
| DELETE | `/bookmarks/{id}` | Soft-delete a bookmark |

## Database Design

### Bookmarks

The bookmarks table stores bookmark information, including:

- `id`
- `title`
- `url`
- `created_at`
- `updated_at`
- `deleted_at`

Soft deletion uses the `deleted_at` timestamp instead of permanently removing a row.

An active bookmark has `deleted_at IS NULL`.

### Users

The users table has been created as part of the authentication stage.

It contains:

- `id`
- `name`
- `email`
- `password_hash`
- `created_at`

Registration and login functionality are still under development.

## Database Migrations

Migrations are maintained in `db/migrations/`.

| Migration | Description |
|---|---|
| `001_create_bookmarks` | Creates the bookmarks table |
| `002_add_bookmark_deleted_at` | Adds soft-delete support |
| `003_create_users` | Creates the users table |

Each migration has an UP and DOWN SQL file.

## Project Structure

```text
Linkboard/
├── cmd/
│   └── server/
│       └── main.go
├── db/
│   ├── migrations/
│   │   ├── 001_create_bookmarks.up.sql
│   │   ├── 001_create_bookmarks.down.sql
│   │   ├── 002_add_bookmark_deleted_at.up.sql
│   │   ├── 002_add_bookmark_deleted_at.down.sql
│   │   ├── 003_create_users.up.sql
│   │   └── 003_create_users.down.sql
│   └── schema.sql
├── internal/
│   ├── database/
│   │   └── postgres.go
│   ├── handler/
│   │   └── bookmark.go
│   ├── repository/
│   │   ├── bookmark.go
│   │   └── user.go
│   └── service/
│       └── bookmark.go
├── frontend/
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Running Locally

### Prerequisites

- Go
- Docker
- Docker Compose
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Ganesh-12-spec/Linkboard.git
cd Linkboard
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

### 3. Configure the database connection

Set the `DATABASE_URL` environment variable according to your local PostgreSQL configuration.

Example:

```bash
export DATABASE_URL="postgres://linkboard:YOUR_PASSWORD@localhost:5433/linkboard"
```

Use your actual local database credentials.

### 4. Start the server

```bash
go run ./cmd/server
```

The server listens on:

```text
http://localhost:8080
```

## Development Roadmap

### Stage A — Backend Foundations
- [x] Go HTTP server
- [x] PostgreSQL connection pooling
- [x] Project architecture

### Stage B — Bookmark CRUD
- [x] Create bookmarks
- [x] List bookmarks
- [x] Get bookmark by ID
- [x] Pagination

### Stage C — Advanced Bookmark Operations
- [x] Update bookmarks
- [x] Soft deletion
- [x] Database migrations
- [ ] Complete practical testing and revision

### Stage D — Authentication and Authorization
- [x] Users table migration
- [x] Initial User model
- [ ] Password hashing with bcrypt
- [ ] User registration
- [ ] User login
- [ ] JWT creation and verification
- [ ] Authentication middleware
- [ ] Bookmark ownership checks

### Future Improvements
- Automated unit and integration tests
- Redis caching
- Search using PostgreSQL full-text search
- Tags and bookmark categorization
- Dockerized application deployment
- CI/CD pipeline
- Logging and observability

## Engineering Goals

The objective of Linkboard is not just to build a working API, but to understand the engineering concepts behind it:

- REST API design
- Layered backend architecture
- SQL and relational database design
- Database migrations
- Authentication and authorization
- Password security
- JWT-based authentication
- Error handling and testing
- Docker and deployment workflows

**Status:** Active development — Bookmark operations implemented; authentication under development.

---

**Author:** Ganesh Mhetre  
**GitHub:** [Ganesh-12-spec](https://github.com/Ganesh-12-spec)
