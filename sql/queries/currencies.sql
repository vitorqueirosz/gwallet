-- name: CreateCurrencies :one
INSERT INTO currencies (id, name, code, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCurrencies :many
SELECT * FROM currencies;

-- name: GetCurrencyByCode :one
SELECT * FROM currencies
WHERE code = $1;


-- name: UpdateCurrencyPrice :many
UPDATE currencies
SET price = $2
WHERE id = $1
RETURNING *;
