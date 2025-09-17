# Go Gin Album API

A simple RESTful web service built with Go and the [Gin framework](https://github.com/gin-gonic/gin). This application provides an API to manage a collection of music albums.

This project is based on the official Go tutorial: [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin).

## Prerequisites

*   Go (version 1.16 or later)

## Getting Started

1.  **Clone the repository:**
    ```sh
    git clone <your-repository-url>
    cd web-service-gin
    ```

2.  **Install dependencies:**
    This will download the Gin framework package.
    ```sh
    go mod tidy
    ```

3.  **Run the server:**
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
        "id": "1",
        "title": "Blue Train",
        "artist": "John Coltrane",
        "price": 56.99
    },
    {
        "id": "2",
        "title": "Jeru",
        "artist": "Gerry Mulligan",
        "price": 17.99
    },
    {
        "id": "3",
        "title": "Sarah Vaughan and Clifford Brown",
        "artist": "Sarah Vaughan",
        "price": 39.99
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
curl http://localhost:8080/albums/2
```

**Example Response (`200 OK`):**
```json
{
    "id": "2",
    "title": "Jeru",
    "artist": "Gerry Mulligan",
    "price": 17.99
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
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```

**Example Response (`201 Created`):**
```json
{
    "id": "4",
    "title": "The Modern Sound of Betty Carter",
    "artist": "Betty Carter",
    "price": 49.99
}
```

