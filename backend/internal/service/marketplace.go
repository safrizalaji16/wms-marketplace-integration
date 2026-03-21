package service

import (
	"backend/domain"
	"backend/dto"
	"context"
)

type marketplaceService struct {
	marketplaceRepository domain.MarketplaceRepository
}

func NewMarketplace(marketplaceRepository domain.MarketplaceRepository) domain.MarketplaceService {
	return &marketplaceService{
		marketplaceRepository: marketplaceRepository,
	}
}

func (m *marketplaceService) Bootstrap(ctx context.Context, state string) (dto.MarketplaceTokenData, error) {
	return m.marketplaceRepository.Bootstrap(ctx, state)
}

func (m *marketplaceService) Authorize(ctx context.Context, state string) (dto.MarketplaceAuthorizeData, error) {
	return m.marketplaceRepository.Authorize(ctx, state)
}

func (m *marketplaceService) ExchangeToken(ctx context.Context, code string) (dto.MarketplaceTokenData, error) {
	return m.marketplaceRepository.ExchangeToken(ctx, code)
}

func (m *marketplaceService) RefreshToken(ctx context.Context, refreshToken string) (dto.MarketplaceTokenData, error) {
	return m.marketplaceRepository.RefreshToken(ctx, refreshToken)
}

func (m *marketplaceService) GetLogisticChannels(ctx context.Context) ([]dto.MarketplaceLogisticChannelData, error) {
	return m.marketplaceRepository.GetLogisticChannels(ctx)
}
