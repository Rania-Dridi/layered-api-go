# layered-api-go

Small REST API in Go showing the layered architecture I used during my final-year project
(Router → Handler → Service → Repository), applied to a simplified course-enrolment domain.

This is an independent demo written for this portfolio. It does not contain any code or data from
the Wledi platform, which belongs to Audasoft Technology.

## Layers

| Layer | Package | Role |
|-------|---------|------|
| Router | `internal/handler` (`Router`) | Maps routes to handlers |
| Handler | `internal/handler` | HTTP in/out, status codes |
| Service | `internal/service` | Business rules: capacity, duplicates, grade range |
| Repository | `internal/repository` | Data access behind an interface (in-memory here, PostgreSQL in production) |

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/courses` | List courses |
| POST | `/api/enrollments` | Enrol a student: `{"courseId":1,"studentId":10}` |
| PUT | `/api/enrollments/{id}/grade` | Set a grade (0–20): `{"grade":15.5}` |

Errors map to 404 (not found), 409 (course full / already enrolled), 422 (invalid grade), 400 (bad JSON).

## Run and test

```bash
go run ./cmd/server
go test ./...
```

Requires Go 1.22+ (uses the standard library only).
