package repository

import (
	"backend/domain"
	"context"

	"github.com/uptrace/bun"
)

type OrderItemRepository struct {
	db *bun.DB
}

func NewOrderItem(db *bun.DB) domain.OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (o *OrderItemRepository) GetByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	var orderItems []domain.OrderItem

	err := o.db.NewSelect().
		Model(&orderItems).
		Where("order_id = ?", orderID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return orderItems, nil
}

func (o *OrderItemRepository) Create(ctx context.Context, orderItem *domain.OrderItem) (domain.OrderItem, error) {
	err := o.db.NewInsert().Model(orderItem).Returning("*").Scan(ctx)

	if err != nil {
		return domain.OrderItem{}, err
	}

	return *orderItem, nil
}
