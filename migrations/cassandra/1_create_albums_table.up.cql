CREATE TABLE IF NOT EXISTS albums (
  id UUID PRIMARY KEY,
  title text,
  artist text,
  price double,
  year int,
  image_url text,
  genre text
);