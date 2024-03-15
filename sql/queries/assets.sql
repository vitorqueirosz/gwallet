-- name: CreateUserAsset :one
INSERT INTO assets (id, currency_id, user_id, amount, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserAssets :many
SELECT assets.id AS asset_id, assets.currency_id AS asset_currency_id, assets.user_id, assets.amount, currencies.id AS currency_id, currencies.code AS currency_code, currencies.price AS currency_price from assets
JOIN users ON assets.user_id = users.id
JOIN currencies ON assets.currency_id = currencies.id
WHERE users.id = $1;