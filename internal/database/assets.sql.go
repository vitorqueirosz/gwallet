// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: assets.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUserAsset = `-- name: CreateUserAsset :one
INSERT INTO assets (id, currency_id, user_id, amount, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, currency_id, user_id, amount, created_at, updated_at
`

type CreateUserAssetParams struct {
	ID         uuid.UUID
	CurrencyID uuid.UUID
	UserID     uuid.UUID
	Amount     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) CreateUserAsset(ctx context.Context, arg CreateUserAssetParams) (Asset, error) {
	row := q.db.QueryRowContext(ctx, createUserAsset,
		arg.ID,
		arg.CurrencyID,
		arg.UserID,
		arg.Amount,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Asset
	err := row.Scan(
		&i.ID,
		&i.CurrencyID,
		&i.UserID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserAssets = `-- name: GetUserAssets :many
SELECT assets.id AS asset_id, assets.currency_id AS asset_currency_id, assets.user_id, assets.amount, currencies.id AS currency_id, currencies.code AS currency_code, currencies.price AS currency_price from assets
JOIN users ON assets.user_id = users.id
JOIN currencies ON assets.currency_id = currencies.id
WHERE users.id = $1
`

type GetUserAssetsRow struct {
	AssetID         uuid.UUID
	AssetCurrencyID uuid.UUID
	UserID          uuid.UUID
	Amount          string
	CurrencyID      uuid.UUID
	CurrencyCode    string
	CurrencyPrice   string
}

func (q *Queries) GetUserAssets(ctx context.Context, id uuid.UUID) ([]GetUserAssetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserAssets, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserAssetsRow
	for rows.Next() {
		var i GetUserAssetsRow
		if err := rows.Scan(
			&i.AssetID,
			&i.AssetCurrencyID,
			&i.UserID,
			&i.Amount,
			&i.CurrencyID,
			&i.CurrencyCode,
			&i.CurrencyPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}