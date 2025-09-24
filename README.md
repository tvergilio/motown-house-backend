# Motown House API Backend

A simple RESTful web service built with Go and the [Gin framework](https://github.com/gin-gonic/gin). This application provides an API to manage a collection of music albums.

This project is based on the official Go tutorial: [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin).

The frontend for this API is a Next.js project, available at: [https://github.com/tvergilio/motown-house](https://github.com/tvergilio/motown-house)

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
        "title": "What's Going On",
        "artist": "Marvin Gaye",
        "price": 39.99,
        "year": 1971,
        "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music112/v4/76/36/2d/76362d74-cb7a-8ef9-104e-cde1d858e9a9/20UMGIM95279.rgb.jpg/100x100bb.jpg",
        "genre": "R&B/Soul"
    },
    {
        "id": 102,
        "title": "Songs in the Key of Life",
        "artist": "Stevie Wonder",
        "price": 42.50,
        "year": 1976,
        "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music118/v4/eb/1f/12/eb1f12ec-474c-63aa-43af-09282f423b9d/00602537004737.rgb.jpg/100x100bb.jpg",
        "genre": "Motown"
    },
    {
        "id": 103,
        "title": "Diana",
        "artist": "Diana Ross",
        "price": 28.75,
        "year": 1980,
        "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/aa/87/1c/aa871c20-95be-38bd-97e3-ecfeb8ec404b/15UMGIM06551.rgb.jpg/100x100bb.jpg",
        "genre": "R&B/Soul"
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
    "id": 101,
    "title": "What's Going On",
    "artist": "Marvin Gaye",
    "price": 39.99,
    "year": 1971,
    "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music112/v4/76/36/2d/76362d74-cb7a-8ef9-104e-cde1d858e9a9/20UMGIM95279.rgb.jpg/100x100bb.jpg",
    "genre": "R&B/Soul"
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
    --data '{"title": "Where Did Our Love Go", "artist": "The Supremes", "price": 9.99, "year": 1964, "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music123/v4/5d/c2/4d/5dc24de8-15d7-16e0-7585-72a2bcc721de/14UMGIM62198.rgb.jpg/100x100bb.jpg", "genre": "R&B/Soul"}'
```

**Example Response (`201 Created`):**
```json
{
    "id": 104,
    "title": "Where Did Our Love Go",
    "artist": "The Supremes",
    "price": 9.99,
    "year": 1964,
    "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music123/v4/5d/c2/4d/5dc24de8-15d7-16e0-7585-72a2bcc721de/14UMGIM62198.rgb.jpg/100x100bb.jpg",
    "genre": "R&B/Soul"
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
    --data '{"title": "Thriller 25", "artist": "Michael Jackson", "price": 45.99, "year": 1982, "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", "genre": "R&B/Soul"}'
```

**Example Response (`200 OK`):**
```json
{
    "id": 101,
    "title": "Thriller 25",
    "artist": "Michael Jackson",
    "price": 45.99,
    "year": 1982,
    "imageUrl": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg",
    "genre": "R&B/Soul"
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
