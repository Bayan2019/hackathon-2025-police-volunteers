-- +goose Up
INSERT INTO users(name, password_hash, iin, phone)
    VALUES ('admin', '$2a$10$vqjSWa5BZEvN/ef7a5pTTOZLImNdubTNVqmyU.7ctiG3kEXyrGk/C', '000000000000', '+0(000)000-00-00');

-- +goose Down
DELETE FROM users;