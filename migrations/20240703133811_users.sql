-- +goose Up
CREATE TABLE users (
    login VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
