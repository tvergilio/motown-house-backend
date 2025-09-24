CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    price NUMERIC NOT NULL,
    year INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    genre TEXT NOT NULL
);