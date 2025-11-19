CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE
);