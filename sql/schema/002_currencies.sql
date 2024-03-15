-- +goose Up

CREATE TABLE currencies (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    price DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE currencies;