package api

import (
	"backend/domain"
	"backend/dto"
	"backend/internal/util"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type marketplaceAPI struct {
	marketplaceService domain.MarketplaceService
}

func NewMarketplace(app *fiber.App, marketplaceService domain.MarketplaceService) {
	m := marketplaceAPI{
		marketplaceService: marketplaceService,
	}

	app.Post("/api/marketplace/oauth/bootstrap", m.Bootstrap)
	app.Get("/api/marketplace/oauth/authorize", m.Authorize)
	app.Post("/api/marketplace/oauth/token", m.ExchangeToken)
	app.Post("/api/marketplace/oauth/refresh", m.RefreshToken)
	app.Get("/api/marketplace/logistic/channels", m.GetLogisticChannels)
}

func (m *marketplaceAPI) Bootstrap(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	res, err := m.marketplaceService.Bootstrap(c, ctx.Query("state"))
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (m *marketplaceAPI) Authorize(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	res, err := m.marketplaceService.Authorize(c, ctx.Query("state"))
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (m *marketplaceAPI) ExchangeToken(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	var req dto.MarketplaceTokenRequest
	if err := ctx.BodyParser(&req); err != nil && len(ctx.Body()) > 0 {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	res, err := m.marketplaceService.ExchangeToken(c, req.Code)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (m *marketplaceAPI) RefreshToken(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	var req dto.MarketplaceRefreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil && len(ctx.Body()) > 0 {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	res, err := m.marketplaceService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (m *marketplaceAPI) GetLogisticChannels(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	res, err := m.marketplaceService.GetLogisticChannels(c)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}
