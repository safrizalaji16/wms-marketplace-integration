package domain

import (
	"context"
	"database/sql"
	"time"
)

type MarketplaceToken struct {
	ID           string         `bun:"id,pk,type:uuid,default:gen_random_uuid()" db:"id"`
	ShopID       string         `bun:"shop_id,notnull,unique" db:"shop_id"`
	AuthCode     sql.NullString `bun:"auth_code" db:"auth_code"`
	AccessToken  sql.NullString `bun:"access_token" db:"access_token"`
	RefreshToken sql.NullString `bun:"refresh_token" db:"refresh_token"`
	ExpiresAt    sql.NullTime   `bun:"expires_at" db:"expires_at"`
	CreatedAt    time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp" db:"created_at"`
	UpdatedAt    time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp" db:"updated_at"`
}

type MarketplaceTokenRepository interface {
	GetByShopID(ctx context.Context, shopID string) (MarketplaceToken, error)
	Upsert(ctx context.Context, token *MarketplaceToken) (MarketplaceToken, error)
}
