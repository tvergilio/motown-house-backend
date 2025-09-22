# Go Gin Album API

A simple RESTful web service built with Go and the [Gin framework](https://github.com/gin-gonic/gin). This application provides an API to manage a collection of music albums.

This project is based on the official Go tutorial: [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin).

## Prerequisites

*   Go (version 1.16 or later)
*   Docker (for running Postgres)

## Getting Started

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/tvergilio/web-service-gin
    cd web-service-gin
    ```

2.  **Install dependencies:**
    This will download the Gin framework package.
    ```sh
    go mod tidy
    ```

3.  **Set up environment variables:**
    -   Copy the provided `.env` file or create one in the project root:
        ```sh
        cp .env.example .env
        ```
    -   Edit `.env` to match your desired Postgres credentials and database name.

4.  **Start Postgres using Docker Compose:**
    ```sh
    docker-compose up -d
    ```
    This will start a Postgres container with persistent storage.

5.  **Run database migrations:**
    ```sh
    migrate -path ./migrations -database "postgres://<user>:<password>@localhost:5432/<db>?sslmode=disable" up
    ```
    Replace `<user>`, `<password>`, and `<db>` with the values from your `.env` file.

6.  **Run the server:**
    ```sh
    go run main.go
    ```
    The server will start and listen on `http://localhost:8080`.

## API Endpoints

The API provides the following endpoints for interacting with the album data.

---

### 1. Get All Albums

Retrieves a list of all albums in the collection.

*   **Method:** `GET`
*   **Endpoint:** `/albums`

**Example Request:**
```sh
curl http://localhost:8080/albums
```

**Example Response (`200 OK`):**
```json
[
    {
        "id": "101",
        "title": "Thriller",
        "artist": "Michael Jackson",
        "price": 42.99
    },
    {
        "id": "102",
        "title": "Lady Soul",
        "artist": "Aretha Franklin",
        "price": 35.50
    },
    {
        "id": "103",
        "title": "What's Going On",
        "artist": "Marvin Gaye",
        "price": 39.00
    }
]
```

---

### 2. Get Album by ID

Retrieves a single album by its unique ID.

*   **Method:** `GET`
*   **Endpoint:** `/albums/:id`

**Example Request:**
```sh
curl http://localhost:8080/albums/102
```

**Example Response (`200 OK`):**
```json
{
    "id": "102",
    "title": "Lady Soul",
    "artist": "Aretha Franklin",
    "price": 35.50
}
```

If an album with the specified ID is not found, it will return a `404 Not Found` with a message.

---

### 3. Create a New Album

Adds a new album to the collection.

*   **Method:** `POST`
*   **Endpoint:** `/albums`
*   **Body:** `JSON`

**Example Request:**
```sh
curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "104","title": "Bad","artist": "Michael Jackson","price": 29.99}'
```

**Example Response (`201 Created`):**
```json
{
    "id": "104",
    "title": "Bad",
    "artist": "Michael Jackson",
    "price": 29.99
}
```

## Persistence Layer

This project uses a PostgreSQL database for persistent storage of album data. The database connection is managed using [sqlx](https://github.com/jmoiron/sqlx) and the [lib/pq](https://github.com/lib/pq) Postgres driver. Database migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate) and are located in the `migrations/` directory. All album data is stored and retrieved from the database, ensuring data is not lost between application restarts.
