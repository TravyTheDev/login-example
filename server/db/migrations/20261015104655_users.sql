-- +goose Up
CREATE TABLE users (
    id integer primary key autoincrement,
    username varchar(255),
    email varchar(255),
    password varchar(255),
    role tinyint,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email)
);

-- +goose Down
DROP TABLE users;
