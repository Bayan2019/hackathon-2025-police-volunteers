-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    revoked_at TEXT
);

-- +goose Down
DROP TABLE refresh_tokens;