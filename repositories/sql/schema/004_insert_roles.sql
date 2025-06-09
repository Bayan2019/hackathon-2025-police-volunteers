-- +goose Up
INSERT INTO roles(title)
VALUES ('volunteer'),
       ('police'),
       ('admin');

-- +goose Down
DELETE FROM roles;