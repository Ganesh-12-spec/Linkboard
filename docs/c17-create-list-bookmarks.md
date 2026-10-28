# C17 — Create and List Bookmarks

## 1. Understand the architecture

**Request flow:**

`Client → Handler → Service → Repository → PostgreSQL`

**Response flow:**

`PostgreSQL → Repository → Service → Handler → Client`

Each layer has one main job. Keeping these jobs separate makes the backend easier to understand, test, and maintain.

## 2. What does each layer do?

### Handler — `internal/handler/bookmark.go`

**Job:** Handle HTTP requests and responses.

- Reads incoming JSON.
- Validates the request data.
- Calls the service layer.
- Returns the correct HTTP status and JSON response.

**Remember:** Handler understands HTTP. It should not execute SQL directly.

### Service — `internal/service/bookmark.go`

**Job:** Coordinate the operation.

- Receives data from the handler.
- Calls the repository.
- Returns the result or error.
- Can contain business rules as the project grows.

**Remember:** Our current service is intentionally simple. It mainly passes the request to the repository and returns the result.

### Repository — `internal/repository/bookmark.go`

**Job:** Communicate with PostgreSQL.

- Executes SQL queries.
- Passes values using SQL parameters.
- Reads database results into Go structs.
- Returns records or errors to the service.

**Remember:** Repository handles database operations, not HTTP status codes.

### Database — PostgreSQL

**Job:** Store bookmark data permanently.

- Saves records in the `bookmarks` table.
- Generates IDs according to the table schema.
- Returns requested columns when a query uses `RETURNING`.

### Database connection — `internal/database/postgres.go`

**Job:** Establish and manage database connections.

- Creates a `pgxpool` connection pool.
- Checks whether the database is reachable.
- Allows the application to reuse database connections.

**Remember:** `pgxpool` manages connections; it does not validate bookmark input.

## 3. C17 — Create a bookmark

Endpoint: `POST /bookmarks`

1. **Handler:** Decode the JSON body and validate the title and URL.
2. **Service:** Call the repository with the request context and bookmark data.
3. **Repository:** Execute an `INSERT` query using SQL parameters.
4. **PostgreSQL:** Save the bookmark and generate its ID.
5. **Repository:** Use `RETURNING` to retrieve the saved fields and `Scan` to copy them into Go values.
6. **Service:** Return the bookmark or error.
7. **Handler:** Send `201 Created` and the JSON response when successful.

### Important SQL concepts

- `INSERT`: adds a new database record.
- `$1`, `$2`: placeholders for parameterized values such as title and URL.
- `RETURNING`: retrieves selected fields from the inserted record.
- `QueryRow`: executes a query expected to return one row.
- `Scan`: copies the returned column values into Go variables.

**Example request:**

`POST /bookmarks`

```json
{
  "title": "Go Documentation",
  "url": "https://go.dev"
}
```

A successful response contains the saved bookmark, including its generated ID, with HTTP status `201 Created`.

## 4. C17 — List bookmarks

Endpoint: `GET /bookmarks`

1. **Handler:** Receives the request and calls the service.
2. **Service:** Calls the repository.
3. **Repository:** Executes a `SELECT` query.
4. **Repository:** Uses `rows.Next()` to move through the result rows.
5. **Repository:** Uses `Scan` to read each row and `append` to collect bookmarks in a slice.
6. **Repository:** Checks `rows.Err()` for errors encountered during iteration.
7. **Service:** Returns the list or error.
8. **Handler:** Sends the bookmarks as JSON.

### Important Go concepts

- `context.Context`: carries cancellation and deadline information through the request.
- `error`: communicates whether an operation failed.
- `make`: creates a slice with a specified length and capacity.
- `append`: adds elements to a slice.
- `rows.Next()`: advances to the next result row.
- `rows.Err()`: reports errors encountered while iterating through rows.
- Dependency injection: the service receives its repository as a dependency instead of constructing the database layer itself.

## 5. HTTP status codes to remember

| Status | Meaning | Example |
|---|---|---|
| `201 Created` | Resource created | Bookmark saved successfully |
| `400 Bad Request` | Invalid request | Empty title or invalid URL |
| `500 Internal Server Error` | Unexpected server failure | Database operation fails unexpectedly |

Invalid input should be rejected before the service attempts to create the bookmark.

## 6. Debugging checklist

If the endpoint fails, check these in order:

1. Is PostgreSQL running?
2. Is `DATABASE_URL` correct?
3. Was the connection pool initialized successfully?
4. Is the route registered for the correct HTTP method?
5. Is the JSON valid and does it match the request struct?
6. Does the SQL match the database schema?
7. Are SQL parameters in the correct order?
8. Does `Scan` match the columns returned by the query?
9. Are errors propagated and handled correctly?
10. Do `curl` and a database query confirm the result?

## 7. What to remember for C18 and C19

### C18 — Fetch one bookmark by ID

The architecture stays the same. The repository now needs a lookup query such as `SELECT ... WHERE id = $1`.

Remember to validate the ID and handle the case where the bookmark does not exist. Distinguish a missing record from an unexpected database error.

### C19 — If implementing update-by-ID

The handler validates the ID and updated data, the service coordinates the operation, and the repository executes a parameterized `UPDATE`.

Remember to handle a nonexistent ID and decide how the endpoint reports the updated record or failure.

## 8. Final revision rule

**Handler = HTTP. Service = coordination and business rules. Repository = SQL. PostgreSQL = persistent storage.**

Before moving to the next endpoint, explain the request flow from memory, rebuild the relevant code independently, and test both successful and invalid requests.