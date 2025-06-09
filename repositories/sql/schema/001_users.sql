-- +goose Up
CREATE TABLE users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    -- email TEXT NOT NULL UNIQUE,
    iin TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL, 
    date_of_birth TEXT NOT NULL DEFAULT '2000-01-01',
    password_hash TEXT NOT NULL,
    current_location TEXT NOT NULL DEFAULT 'Almaty'
);

-- +goose Down
DROP TABLE users;