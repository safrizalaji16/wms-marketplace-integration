package service

import (
	"backend/domain"
	"backend/dto"
	"backend/internal/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
)

type OrderService struct {
	OrderRepository       domain.OrderRepository
	OrderItemRepository   domain.OrderItemRepository
	MarketplaceRepository domain.MarketplaceRepository
}

func NewOrderService(orderRepository domain.OrderRepository, orderItemRepository domain.OrderItemRepository, marketplaceRepository domain.MarketplaceRepository) domain.OrderService {
	return &OrderService{
		OrderRepository:       orderRepository,
		OrderItemRepository:   orderItemRepository,
		MarketplaceRepository: marketplaceRepository,
	}
}

func (o *OrderService) GetOrders(ctx context.Context, filter dto.OrderFilter) (dto.OrderListData, error) {
	totalOrders, err := o.OrderRepository.Count(ctx)
	if err != nil {
		return dto.OrderListData{}, err
	}

	if totalOrders == 0 {
		if err := o.syncOrders(ctx); err != nil {
			return dto.OrderListData{}, err
		}
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	filter.Offset = (filter.Page - 1) * filter.Limit

	orders, total, err := o.OrderRepository.GetAll(ctx, filter)

	if err != nil {
		return dto.OrderListData{}, err
	}

	var orderData []dto.OrderData
	for _, order := range orders {
		orderData = append(orderData, dto.OrderData{
			ID:                order.ID,
			OrderSN:           order.OrderSN,
			WMSStatus:         order.WMSStatus,
			MarketplaceStatus: order.MarketplaceStatus,
			ShippingStatus:    order.ShippingStatus,
			TrackingNumber:    order.TrackingNumber,
			UpdatedAt:         util.NullableTime(order.UpdatedAt),
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(filter.Limit)))
	}

	return dto.OrderListData{
		Orders:     orderData,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (o *OrderService) syncOrders(ctx context.Context) error {
	if o.MarketplaceRepository == nil {
		return util.ErrMarketplaceNotConfigured
	}

	marketplaceOrders, err := o.MarketplaceRepository.GetOrders(ctx)
	if err != nil {
		slog.Error("sync orders failed while fetching marketplace orders", "error", err.Error())
		return err
	}

	slog.Info("sync orders started", "marketplace_order_count", len(marketplaceOrders))

	for _, marketplaceOrder := range marketplaceOrders {
		if !isEligibleForWMSSync(marketplaceOrder) {
			slog.Info("sync order skipped: not eligible", "order_sn", marketplaceOrder.OrderSN, "marketplace_status", marketplaceOrder.MarketplaceStatus, "shipping_status", marketplaceOrder.ShippingStatus)
			continue
		}

		existingOrder, err := o.OrderRepository.GetByOrderSN(ctx, marketplaceOrder.OrderSN)
		if err == nil && existingOrder.ID != "" {
			slog.Info("sync order skipped: already exists", "order_sn", marketplaceOrder.OrderSN)
			continue
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			slog.Error("sync orders failed while checking existing order", "order_sn", marketplaceOrder.OrderSN, "error", err.Error())
			return err
		}
		if err := o.createOrderFromMarketplace(ctx, marketplaceOrder); err != nil {
			return err
		}
	}

	return nil
}

func isEligibleForWMSSync(order domain.Order) bool {
	if strings.ToLower(strings.TrimSpace(order.MarketplaceStatus)) != "paid" {
		return false
	}

	// switch strings.ToLower(strings.TrimSpace(order.ShippingStatus)) {
	// case "shipped", "delivered", "cancelled":
	// 	return false
	// default:
	// 	return true
	// }

	return true
}

func (o *OrderService) GetByOrderSN(ctx context.Context, orderSN string) (dto.OrderDetailData, error) {
	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)

	if err != nil {
		return dto.OrderDetailData{}, err
	}

	if order.ID == "" {
		return dto.OrderDetailData{}, util.ErrOrderNotFound
	}

	if o.OrderItemRepository == nil {
		return dto.OrderDetailData{}, errors.New("order item repository is not configured")
	}

	orderItems, err := o.OrderItemRepository.GetByOrderID(ctx, order.ID)
	if err != nil {
		return dto.OrderDetailData{}, err
	}

	var items []dto.OrderItemData
	for _, item := range orderItems {
		items = append(items, dto.OrderItemData{
			ID:       item.ID,
			OrderID:  item.OrderID,
			SKU:      item.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	return dto.OrderDetailData{
		ID:                order.ID,
		OrderSN:           order.OrderSN,
		ShopID:            order.ShopID,
		WMSStatus:         order.WMSStatus,
		MarketplaceStatus: order.MarketplaceStatus,
		ShippingStatus:    order.ShippingStatus,
		TrackingNumber:    order.TrackingNumber,
		TotalAmount:       order.TotalAmount,
		Items:             items,
	}, nil
}

func (o *OrderService) Pick(ctx context.Context, orderSN string) error {
	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.ID == "" {
		return util.ErrOrderNotFound
	}

	if order.WMSStatus != "READY_TO_PICK" {
		return fmt.Errorf("%w: order cannot be picked", util.ErrInvalidStatusTransition)
	}

	order.WMSStatus = "PICKING"
	slog.Info("order picked", "order_sn", orderSN, "next_wms_status", order.WMSStatus)
	return o.OrderRepository.Update(ctx, &order)
}

func (o *OrderService) Pack(ctx context.Context, orderSN string) error {
	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.ID == "" {
		return util.ErrOrderNotFound
	}

	if order.WMSStatus != "PICKING" {
		return fmt.Errorf("%w: order cannot be packed", util.ErrInvalidStatusTransition)
	}

	order.WMSStatus = "PACKED"
	slog.Info("order packed", "order_sn", orderSN, "next_wms_status", order.WMSStatus)
	return o.OrderRepository.Update(ctx, &order)
}

func (o *OrderService) Ship(ctx context.Context, orderSN, channelID string) error {
	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.ID == "" {
		return util.ErrOrderNotFound
	}

	if order.WMSStatus != "PACKED" {
		return fmt.Errorf("%w: order cannot be shipped", util.ErrInvalidStatusTransition)
	}

	if o.MarketplaceRepository == nil {
		return util.ErrMarketplaceNotConfigured
	}
	if strings.TrimSpace(channelID) == "" {
		return util.ErrChannelIDRequired
	}

	shipData, err := o.MarketplaceRepository.ShipOrder(ctx, orderSN, channelID)
	if err != nil {
		return err
	}

	order.WMSStatus = "SHIPPED"
	order.ShippingStatus = shipData.ShippingStatus
	order.TrackingNumber = shipData.TrackingNo
	slog.Info("order shipped", "order_sn", orderSN, "channel_id", channelID, "shipping_status", shipData.ShippingStatus, "tracking_number", shipData.TrackingNo)

	return o.OrderRepository.Update(ctx, &order)
}

func (o *OrderService) UpdateMarketplaceStatus(ctx context.Context, orderSN, status string) error {
	orderSN = strings.TrimSpace(orderSN)
	status = util.NormalizeStatus(status)

	if orderSN == "" {
		return util.ErrOrderSNRequired
	}
	if status == "" {
		return util.ErrStatusRequired
	}
	if !util.IsSupportedMarketplaceWebhookStatus(status) {
		return util.ErrUnsupportedMarketplaceStatus
	}

	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if status == "paid" {
				return o.syncPaidOrderFromMarketplace(ctx, orderSN)
			}
			return util.ErrOrderNotFound
		}
		return err
	}

	if order.ID == "" {
		if status == "paid" {
			return o.syncPaidOrderFromMarketplace(ctx, orderSN)
		}
		return util.ErrOrderNotFound
	}

	currentStatus := util.NormalizeStatus(order.MarketplaceStatus)
	if util.IsTerminalMarketplaceStatus(currentStatus) && currentStatus != status {
		slog.Info("order status webhook ignored: terminal status", "order_sn", orderSN, "current_status", currentStatus, "incoming_status", status)
		return nil
	}

	order.MarketplaceStatus = status
	slog.Info("order status webhook accepted", "order_sn", orderSN, "status", status)
	return o.OrderRepository.Update(ctx, &order)
}

func (o *OrderService) syncPaidOrderFromMarketplace(ctx context.Context, orderSN string) error {
	if o.MarketplaceRepository == nil {
		return util.ErrMarketplaceNotConfigured
	}

	marketplaceOrder, err := o.MarketplaceRepository.GetOrderBySN(ctx, orderSN)
	if err != nil {
		slog.Error("paid webhook sync failed while fetching marketplace order", "order_sn", orderSN, "error", err.Error())
		return err
	}

	if !isEligibleForWMSSync(marketplaceOrder) {
		slog.Info("paid webhook sync skipped: order not eligible", "order_sn", orderSN, "marketplace_status", marketplaceOrder.MarketplaceStatus, "shipping_status", marketplaceOrder.ShippingStatus)
		return nil
	}

	return o.createOrderFromMarketplace(ctx, marketplaceOrder)
}

func (o *OrderService) createOrderFromMarketplace(ctx context.Context, marketplaceOrder domain.Order) error {
	order := domain.Order{
		OrderSN:               marketplaceOrder.OrderSN,
		ShopID:                marketplaceOrder.ShopID,
		MarketplaceStatus:     marketplaceOrder.MarketplaceStatus,
		ShippingStatus:        marketplaceOrder.ShippingStatus,
		WMSStatus:             "READY_TO_PICK",
		TrackingNumber:        marketplaceOrder.TrackingNumber,
		TotalAmount:           marketplaceOrder.TotalAmount,
		RawMarketplacePayload: marketplaceOrder.RawMarketplacePayload,
		CreatedAt:             marketplaceOrder.CreatedAt,
	}

	if _, err := o.OrderRepository.Create(ctx, &order); err != nil {
		slog.Error("sync order failed while creating order", "order_sn", marketplaceOrder.OrderSN, "error", err.Error())
		return err
	}

	slog.Info("sync order created", "order_sn", order.OrderSN, "item_count", len(marketplaceOrder.Items))

	if o.OrderItemRepository == nil {
		return errors.New("order item repository is not configured")
	}

	for _, item := range marketplaceOrder.Items {
		orderItem := domain.OrderItem{
			OrderID:  order.ID,
			SKU:      item.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
		}

		if _, err := o.OrderItemRepository.Create(ctx, &orderItem); err != nil {
			slog.Error("sync order item failed", "order_sn", marketplaceOrder.OrderSN, "sku", item.SKU, "error", err.Error())
			return err
		}
	}

	return nil
}

func (o *OrderService) UpdateShippingStatus(ctx context.Context, orderSN, status string) error {
	orderSN = strings.TrimSpace(orderSN)
	status = util.NormalizeStatus(status)

	if orderSN == "" {
		return util.ErrOrderSNRequired
	}
	if status == "" {
		return util.ErrStatusRequired
	}
	if !util.IsSupportedShippingWebhookStatus(status) {
		return util.ErrUnsupportedShippingStatus
	}

	order, err := o.OrderRepository.GetByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.ID == "" {
		return util.ErrOrderNotFound
	}

	currentShippingStatus := util.NormalizeStatus(order.ShippingStatus)
	if util.IsTerminalShippingStatus(currentShippingStatus) && currentShippingStatus != status {
		slog.Info("shipping status webhook ignored: terminal status", "order_sn", orderSN, "current_status", currentShippingStatus, "incoming_status", status)
		return nil
	}

	order.ShippingStatus = status
	slog.Info("shipping status webhook accepted", "order_sn", orderSN, "status", status)
	return o.OrderRepository.Update(ctx, &order)
}
