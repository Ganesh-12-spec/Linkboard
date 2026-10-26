# Go Backend Engineering — Database Setup Notes

These notes cover general backend concepts learned while building Linkboard. The same principles apply to many Go backend projects.

## 1. Git: Commits and Pushes

- `git status` — shows changed, staged, and untracked files.
- `git diff` — shows unstaged changes.
- `git add <file>` — stages changes for a commit.
- `git commit -m "message"` — records staged changes in local Git history.
- `git push origin main` — uploads local commits to the remote `main` branch.
- A commit records a real change. Rewriting a file identically does not create a new commit.
- Use small, meaningful commits that represent one logical change.

## 2. Docker Compose and PostgreSQL

Docker Compose defines and runs services using a YAML configuration file.

In Linkboard, the PostgreSQL service specifies the database image, credentials, port mapping, persistent volume, and health check.

Important concepts:

- **Image:** the template used to create a container.
- **Container:** the running instance of an image.
- **Port mapping:** makes a container service accessible through a host port.
- **Volume:** preserves database files beyond the lifetime of a container.
- **Health check:** checks whether the configured service is ready.

A mapping such as `127.0.0.1:5433:5432` means host port `5433` forwards to container port `5432`, accessible through the host's loopback interface.

Development credentials must not be reused as production secrets.

## 3. Database Schema

A schema describes the structure of database objects, such as tables, columns, data types, constraints, and indexes.

For example, a bookmarks table might contain an ID, URL, title, and timestamps.

The schema defines what the database structure should look like. It does not, by itself, explain how to upgrade an existing database safely over time.

## 4. Database Migrations

A migration is a versioned change to a database structure.

Common naming conventions:

- `001_create_bookmarks.up.sql` — applies the change.
- `001_create_bookmarks.down.sql` — attempts to reverse the change.

The `.up.sql` file might create a table; its corresponding down migration might drop that table.

Important principles:

- Apply migrations in a known order.
- Track which migrations have already been applied.
- Review destructive changes carefully.
- A down migration may lose data and is not always a perfect reversal.
- Keep schema changes under version control.

A schema snapshot and migration files serve related but different purposes: a snapshot describes a structure, while migrations describe changes to that structure.

## 5. Go Packages and Imports

A Go source file starts with a package declaration.

- `package main` is used for an executable entry-point package.
- `package database` groups database-related code into a reusable package.
- `import` makes other packages available to the file.
- `func` declares a function.
- `:=` declares and initializes a variable with an inferred type.

Keep related responsibilities organized into packages instead of putting every operation into one large file.

## 6. PostgreSQL Connection Pooling with pgx

`github.com/jackc/pgx/v5/pgxpool` provides PostgreSQL connection-pool functionality for Go.

A connection pool manages reusable database connections. Reusing connections avoids the overhead of establishing a fresh connection for every operation and helps control concurrent database access.

A pool is a connection manager, not the database itself.

Creating a pool does not necessarily prove that the database is reachable. A startup ping is useful for checking connectivity before the application begins serving requests.

## 7. Connection Strings and Environment Variables

A connection string contains information needed to connect to a database, such as the host, port, database name, username, and authentication details.

Go's `os.Getenv("DATABASE_URL")` reads an environment variable.

Checking whether the value is empty verifies that configuration was supplied. It does not verify that the credentials work or that PostgreSQL is reachable.

Keep production secrets out of source code and public repositories.

## 8. Context, Timeouts, and Cancellation

Go's `context` package lets operations carry cancellation signals and deadlines.

- `context.Background()` provides a root context.
- `context.WithTimeout(parent, duration)` creates a context with a time limit, subject to the parent context's earlier deadline or cancellation.
- `cancel()` releases resources associated with the derived context and signals cancellation.
- `defer cancel()` schedules cancellation when the surrounding function returns.

A timeout limits how long an operation is allowed to wait; it does not guarantee success or make the operation faster.

## 9. Ping, Errors, and Cleanup

`pool.Ping(ctx)` checks whether PostgreSQL is reachable using the supplied context.

Go functions can return multiple values, including a result and an error. For example, a pool-initialization function can return `(*pgxpool.Pool, error)`.

Common patterns:

- Return a nil result with an error when initialization fails.
- Use `fmt.Errorf("operation: %w", err)` to add context while wrapping an error.
- Check errors explicitly.
- Close a pool if initialization fails after creating it.
- Use `defer pool.Close()` in the component responsible for the pool's lifetime.

A successful startup ping confirms connectivity at that moment; later database operations can still fail.

## 10. Typical Backend Responsibilities

A common layered backend design is:

- **Handler:** receives HTTP requests and writes HTTP responses.
- **Service:** implements application or business logic.
- **Repository:** performs database queries and persistence operations.
- **Database package:** initializes and manages database connectivity.

The exact design varies by project. A connection pool must be passed to the components that need it; creating the pool alone does not automatically connect it to handlers or repositories.

## 11. What I Have Reviewed in Linkboard

- PostgreSQL development setup through Docker Compose.
- Database schema and versioned migration files.
- Go's PostgreSQL connection-pool initialization.
- Environment-variable validation.
- Context timeouts, connectivity checks, error wrapping, and cleanup.

## 12. Topics to Practise Next

- How thinfrastructure to read and modify data.
