DROP TABLE IF EXISTS albums;
CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    price NUMERIC NOT NULL,
    year INTEGER NOT NULL,
    image_url TEXT,
    genre TEXT
);