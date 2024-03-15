-- +goose Up

CREATE TABLE assets (
    id UUID PRIMARY KEY,
    currency_id UUID UNIQUE NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE assets;