package api

import (
	"backend/domain"
	"backend/dto"
	"backend/internal/middleware"
	"backend/internal/util"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type orderAPI struct {
	orderService domain.OrderService
}

func NewOrder(app *fiber.App, orderService domain.OrderService, auzMidd fiber.Handler) {
	o := orderAPI{
		orderService: orderService,
	}

	app.Post("/webhook/order-status", o.OrderStatusWebhook)
	app.Post("/webhook/shipping-status", o.ShippingStatusWebhook)
	app.Use(auzMidd)
	app.Get("/api/orders", o.GetOrders)
	app.Get("/api/orders/:order_sn", o.GetOrderByOrderSN)
	app.Post(
		"/api/orders/:order_sn/pick",
		middleware.AuthorizeRoles("picker", "superAdmin"),
		o.PickOrder,
	)
	app.Post(
		"/api/orders/:order_sn/pack",
		middleware.AuthorizeRoles("packer", "superAdmin"),
		o.PackOrder,
	)
	app.Post(
		"/api/orders/:order_sn/ship",
		middleware.AuthorizeRoles("warehouseAdmin", "superAdmin"),
		o.ShipOrder,
	)
}

func (o *orderAPI) OrderStatusWebhook(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.OrderStatusWebhookRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	if err := o.orderService.UpdateMarketplaceStatus(c, req.OrderSN, req.Status); err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.Response[map[string]string]{
		Message: "Order status updated",
		Data: map[string]string{
			"order_sn": req.OrderSN,
			"status":   req.Status,
		},
	})
}

func (o *orderAPI) ShippingStatusWebhook(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.ShippingStatusWebhookRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	if err := o.orderService.UpdateShippingStatus(c, req.OrderSN, req.Status); err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.Response[map[string]string]{
		Message: "Shipping status updated",
		Data: map[string]string{
			"order_sn":       req.OrderSN,
			"shipping_state": req.Status,
		},
	})
}

func (o *orderAPI) GetOrders(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	filter := dto.OrderFilter{
		Search:            ctx.Query("search"),
		WMSStatus:         ctx.Query("wms_status"),
		MarketplaceStatus: ctx.Query("marketplace_status"),
		ShippingStatus:    ctx.Query("shipping_status"),
		Page:              parseIntOrDefault(ctx.Query("page"), 1),
		Limit:             parseIntOrDefault(ctx.Query("limit"), 10),
	}

	res, err := o.orderService.GetOrders(c, filter)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (o *orderAPI) GetOrderByOrderSN(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	orderSN := ctx.Params("order_sn")

	res, err := o.orderService.GetByOrderSN(c, orderSN)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (o *orderAPI) PickOrder(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	orderSN := ctx.Params("order_sn")

	if err := o.orderService.Pick(c, orderSN); err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(map[string]string{
		"order_sn":   orderSN,
		"wms_status": "PICKING",
	}))
}

func (o *orderAPI) PackOrder(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	orderSN := ctx.Params("order_sn")

	if err := o.orderService.Pack(c, orderSN); err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(map[string]string{
		"order_sn":   orderSN,
		"wms_status": "PACKED",
	}))
}

func (o *orderAPI) ShipOrder(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	orderSN := ctx.Params("order_sn")
	var req dto.ShipOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	if err := o.orderService.Ship(c, orderSN, req.ChannelID); err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	order, err := o.orderService.GetByOrderSN(c, orderSN)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(map[string]string{
		"order_sn":        order.OrderSN,
		"wms_status":      order.WMSStatus,
		"shipping_status": order.ShippingStatus,
		"tracking_number": order.TrackingNumber,
	}))
}

func parseIntOrDefault(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
