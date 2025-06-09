-- +goose Up
CREATE TABLE roles(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE roles;