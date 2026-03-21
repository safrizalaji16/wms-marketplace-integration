package repository

import (
	"backend/domain"
	"backend/dto"
	"context"
	"strings"

	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db *bun.DB
}

func NewOrder(db *bun.DB) domain.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o *OrderRepository) Count(ctx context.Context) (int, error) {
	return o.db.NewSelect().Model((*domain.Order)(nil)).Count(ctx)
}

func (o *OrderRepository) GetAll(ctx context.Context, filter dto.OrderFilter) ([]domain.Order, int, error) {
	var orders []domain.Order
	query := o.db.NewSelect().Model(&orders).Order(buildOrderSort(filter))
	countQuery := o.db.NewSelect().Model((*domain.Order)(nil))

	applyOrderFilters(query, filter)
	applyOrderFilters(countQuery, filter)

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query.Offset(filter.Offset)
	}

	err = query.Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func applyOrderFilters(query *bun.SelectQuery, filter dto.OrderFilter) {
	search := strings.TrimSpace(filter.Search)
	wmsStatuses := splitFilterValues(filter.WMSStatus)
	marketplaceStatuses := splitFilterValues(filter.MarketplaceStatus)
	shippingStatuses := splitFilterValues(filter.ShippingStatus)

	if search != "" {
		query.Where("(order_sn ILIKE ? OR tracking_number ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	if len(wmsStatuses) > 0 {
		query.Where("wms_status IN (?)", bun.In(wmsStatuses))
	}

	if len(marketplaceStatuses) > 0 {
		query.Where("marketplace_status IN (?)", bun.In(marketplaceStatuses))
	}

	if len(shippingStatuses) > 0 {
		query.Where("shipping_status IN (?)", bun.In(shippingStatuses))
	}
}

func splitFilterValues(value string) []string {
	rawValues := strings.Split(value, ",")
	values := make([]string, 0, len(rawValues))

	for _, raw := range rawValues {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		values = append(values, trimmed)
	}

	return values
}

func buildOrderSort(filter dto.OrderFilter) string {
	column := normalizeSortColumn(filter.SortBy)
	direction := normalizeSortDirection(filter.SortDir)

	return column + " " + direction
}

func normalizeSortColumn(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "order_sn":
		return "order_sn"
	case "marketplace_status":
		return "marketplace_status"
	case "shipping_status":
		return "shipping_status"
	case "wms_status":
		return "wms_status"
	case "tracking_number":
		return "tracking_number"
	case "updated_at":
		return "updated_at"
	default:
		return "updated_at"
	}
}

func normalizeSortDirection(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "asc":
		return "ASC"
	case "desc":
		return "DESC"
	default:
		return "DESC"
	}
}

func (o *OrderRepository) GetByOrderSN(ctx context.Context, orderSN string) (domain.Order, error) {
	var order domain.Order

	err := o.db.NewSelect().Model(&order).Where("order_sn = ?", orderSN).Scan(ctx)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (o *OrderRepository) Create(ctx context.Context, order *domain.Order) (domain.Order, error) {
	err := o.db.NewInsert().
		Model(order).
		On("CONFLICT (order_sn) DO NOTHING").
		Returning("*").
		Scan(ctx)

	if err != nil {
		return domain.Order{}, err
	}

	return *order, nil
}

func (o *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	_, err := o.db.NewUpdate().
		Model(order).
		Where("order_sn = ?", order.OrderSN).
		Exec(ctx)

	return err
}
