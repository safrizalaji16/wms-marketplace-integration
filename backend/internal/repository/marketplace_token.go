package repository

import (
	"backend/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type MarketplaceTokenRepository struct {
	db *bun.DB
}

func NewMarketplaceToken(db *bun.DB) domain.MarketplaceTokenRepository {
	return &MarketplaceTokenRepository{
		db: db,
	}
}

func (m *MarketplaceTokenRepository) GetByShopID(ctx context.Context, shopID string) (domain.MarketplaceToken, error) {
	var token domain.MarketplaceToken

	err := m.db.NewSelect().
		Model(&token).
		Where("shop_id = ?", shopID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.MarketplaceToken{
			ShopID: shopID,
		}, nil
	}
	if err != nil {
		return domain.MarketplaceToken{}, err
	}

	return token, nil
}

func (m *MarketplaceTokenRepository) Upsert(ctx context.Context, token *domain.MarketplaceToken) (domain.MarketplaceToken, error) {
	existingToken, err := m.GetByShopID(ctx, token.ShopID)
	if err != nil {
		return domain.MarketplaceToken{}, err
	}

	if !token.AuthCode.Valid && existingToken.AuthCode.Valid {
		token.AuthCode = existingToken.AuthCode
	}
	if !token.AccessToken.Valid && existingToken.AccessToken.Valid {
		token.AccessToken = existingToken.AccessToken
	}
	if !token.RefreshToken.Valid && existingToken.RefreshToken.Valid {
		token.RefreshToken = existingToken.RefreshToken
	}
	if !token.ExpiresAt.Valid && existingToken.ExpiresAt.Valid {
		token.ExpiresAt = existingToken.ExpiresAt
	}

	err = m.db.NewInsert().
		Model(token).
		On("CONFLICT (shop_id) DO UPDATE").
		Set("auth_code = EXCLUDED.auth_code").
		Set("access_token = EXCLUDED.access_token").
		Set("refresh_token = EXCLUDED.refresh_token").
		Set("expires_at = EXCLUDED.expires_at").
		Set("updated_at = CURRENT_TIMESTAMP").
		Returning("*").
		Scan(ctx)

	if err != nil {
		return domain.MarketplaceToken{}, err
	}

	return *token, nil
}
