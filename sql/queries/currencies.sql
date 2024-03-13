-- name: CreateCurrencies :one
INSERT INTO currencies (id, name, code, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCurrencies :many
SELECT * FROM currencies;