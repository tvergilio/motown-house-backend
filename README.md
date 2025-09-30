# Motown House API Backend

A RESTful API for managing music albums built with Go and Gin. Features CRUD operations and iTunes integration.

**Frontend**: [Next.js project](https://github.com/tvergilio/motown-house)

## Architecture

This API is built using a clean layered architecture that promotes separation of concerns and maintainability. The system is organised into distinct layers: the API layer manages HTTP requests and responses, the business logic layer handles validation and processing rules, and the repository layer provides a clean abstraction for data access to both PostgreSQL database and external iTunes API.

<img src="diagrams/images/high-level-architecture.svg" width="100%">

### Key Components
- **Repository Pattern**: Clean separation between business logic and data access
- **PostgreSQL**: Primary data storage with migrations
- **iTunes Integration**: External API for album search
- **Comprehensive Testing**: Unit tests with testcontainers for integration testing

## Quick Start

**Prerequisites**: Go 1.25+, Docker, [golang-migrate](https://github.com/golang-migrate/migrate)

```bash
# 1. Clone and setup
git clone https://github.com/tvergilio/web-service-gin
cd web-service-gin
go mod tidy

# 2. Environment (create .env file)
# Use `DB_BACKEND` to switch between `postgres` and `cassandra`.
# Provide a full `POSTGRES_URL` when DB_BACKEND=postgres. 
DB_BACKEND=postgres
POSTGRES_URL=

# Cassandra (only used when DB_BACKEND=cassandra)
# CASSANDRA_HOSTS=host1,host2
# CASSANDRA_KEYSPACE=your_keyspace

# 3. Start services
docker-compose up -d

# 4. Setup keyspace and run migrations
# For Postgres:
migrate -path ./migrations/postgres -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# For Cassandra:
# Step 1: Create keyspace (one-time setup)
docker exec -it cassandra cqlsh -e "CREATE KEYSPACE IF NOT EXISTS motown WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};"
# Step 2: Run table migrations
migrate -path ./migrations/cassandra -database "cassandra://localhost:9042/motown" up

# 5. Start application
go run main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

The following diagrams illustrate how requests are processed through the system:

**Component Interactions**: Shows how components communicate
<img src="diagrams/images/activity.svg" width="100%">

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/albums` | Get all albums |
| GET | `/albums/:id` | Get album by ID |
| POST | `/albums` | Create new album |
| PUT | `/albums/:id` | Update album |
| DELETE | `/albums/:id` | Delete album |
| GET | `/api/search?term=X` | Search iTunes for albums |

**Request Flow**: Shows the logical flow and decision points
<img src="diagrams/images/request-flow.svg" width="100%">

### Example Usage

```bash
# Get all albums
curl http://localhost:8080/albums

# Create album
curl -X POST http://localhost:8080/albums \
  -H "Content-Type: application/json" \
  -d '{"title": "Thriller", "artist": "Michael Jackson", "price": 25.99, "year": 1982, "imageUrl": "...", "genre": "Pop"}'

# Search iTunes
curl "http://localhost:8080/api/search?term=thriller"
```

## Development

```bash
# Run tests (excludes slow Cassandra integration tests)
go test ./...

# Run all tests including Cassandra integration tests (slow, requires Docker)
go test -tags=integration -timeout=300s ./repository

# Docker deployment
docker-compose up --build

# Integration tests only
go test ./repository/...
```

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL with [sqlx](https://github.com/jmoiron/sqlx)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Testing**: [testify](https://github.com/stretchr/testify) + [testcontainers](https://github.com/testcontainers/testcontainers-go)
- **External API**: iTunes Search API integration
