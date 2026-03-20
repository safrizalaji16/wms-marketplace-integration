package domain

import (
	"backend/dto"
	"context"
)

type MarketplaceRepository interface {
	Bootstrap(ctx context.Context, state string) (dto.MarketplaceTokenData, error)
	Authorize(ctx context.Context, state string) (dto.MarketplaceAuthorizeData, error)
	ExchangeToken(ctx context.Context, code string) (dto.MarketplaceTokenData, error)
	RefreshToken(ctx context.Context, refreshToken string) (dto.MarketplaceTokenData, error)
	GetLogisticChannels(ctx context.Context) ([]dto.MarketplaceLogisticChannelData, error)
	ShipOrder(ctx context.Context, orderSN, channelID string) (dto.MarketplaceShipData, error)
	GetOrders(ctx context.Context) ([]Order, error)
	GetOrderBySN(ctx context.Context, orderSN string) (Order, error)
}

type MarketplaceService interface {
	Bootstrap(ctx context.Context, state string) (dto.MarketplaceTokenData, error)
	Authorize(ctx context.Context, state string) (dto.MarketplaceAuthorizeData, error)
	ExchangeToken(ctx context.Context, code string) (dto.MarketplaceTokenData, error)
	RefreshToken(ctx context.Context, refreshToken string) (dto.MarketplaceTokenData, error)
	GetLogisticChannels(ctx context.Context) ([]dto.MarketplaceLogisticChannelData, error)
}
