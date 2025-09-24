# Go Gin Album API

A simple RESTful web service built with Go and the [Gin framework](https://github.com/gin-gonic/gin). This application provides an API to manage a collection of music albums.

This project is based on the official Go tutorial: [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin).

## Prerequisites

* Go (version 1.16 or later)
* Docker (for running Postgres)

## Getting Started

1. **Clone the repository:**
    ```sh
    git clone https://github.com/tvergilio/web-service-gin
    cd web-service-gin
    ```

2. **Install dependencies:**
    This will download the Gin framework package.
    ```sh
    go mod tidy
    ```

3. **Set up environment variables:**
    - Copy or create a `.env` file in the project root:
        ```sh
        cp .env.example .env
        ```
    - Edit `.env` to match your Postgres credentials and database name.

4. **Start Postgres using Docker Compose:**
    ```sh
    docker-compose up -d
    ```

5. **Run database migrations:**
    ```sh
    migrate -path ./migrations -database "postgres://<user>:<password>@localhost:5432/<db>?sslmode=disable" up
    ```
    Replace `<user>`, `<password>`, and `<db>` with your `.env` values.

6. **Run the server:**
    ```sh
    go run main.go
    ```
    The server will listen on `http://localhost:8080`.

## API Endpoints

---

### Get All Albums

* **GET** `/albums`

Returns all albums.

**Example Response (`200 OK`):**
```json
[
    {
        "id": 101,
        "title": "Thriller",
        "artist": "Michael Jackson",
        "price": 42.99,
        "year": 1982
    },
    {
        "id": 102,
        "title": "Lady Soul",
        "artist": "Aretha Franklin",
        "price": 35.50,
        "year": 1968
    },
    {
        "id": 103,
        "title": "What's Going On",
        "artist": "Marvin Gaye",
        "price": 39.00,
        "year": 1971
    }
]
```

---

### Get Album by ID

* **GET** `/albums/:id`

Returns a single album by its ID.

**Example Response (`200 OK`):**
```json
{
    "id": 102,
    "title": "Lady Soul",
    "artist": "Aretha Franklin",
    "price": 35.50,
    "year": 1968
}
```

Returns `404 Not Found` if the album does not exist.

---

### Create a New Album

* **POST** `/albums`
* **Body:** JSON

**Example Request:**
```sh
curl http://localhost:8080/albums \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": 104, "title": "Bad", "artist": "Michael Jackson", "price": 29.99, "year": 1987}'
```

**Example Response (`201 Created`):**
```json
{
    "id": 104,
    "title": "Bad",
    "artist": "Michael Jackson",
    "price": 29.99,
    "year": 1987
}
```

---

### Update an Album

* **PUT** `/albums/:id`
* **Body:** JSON

Updates an album by its ID. The ID in the URI is used for the update.

**Example Request:**
```sh
curl http://localhost:8080/albums/101 \
    --header "Content-Type: application/json" \
    --request "PUT" \
    --data '{"title": "Thriller 25", "artist": "Michael Jackson", "price": 45.99, "year": 1982}'
```

**Example Response (`200 OK`):**
```json
{
    "id": 101,
    "title": "Thriller 25",
    "artist": "Michael Jackson",
    "price": 45.99,
    "year": 1982
}
```

Returns `400 Bad Request` for invalid JSON or ID, and `500 Internal Server Error` if the album does not exist.

---

### Delete an Album

* **DELETE** `/albums/:id`

Deletes an album by its ID.

Returns `204 No Content` on success, or `404 Not Found` if the album does not exist.

## Persistence Layer

Uses PostgreSQL for album storage. Database connection is managed with [sqlx](https://github.com/jmoiron/sqlx) and [lib/pq](https://github.com/lib/pq). Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate) and located in the `migrations/` directory.
