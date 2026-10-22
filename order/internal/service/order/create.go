package order

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"order/internal/metrics"
	"order/internal/model"
	"platform/pkg/logger"
)

func (s *srv) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		metrics.OrdersCreateFailedTotal.WithLabelValues("empty_parts").Inc()
		return nil, model.ErrEmptyPartsList
	}

	qtyByPart := make(map[uuid.UUID]int, len(partUUIDs))
	unique := make([]uuid.UUID, 0, len(partUUIDs))
	for _, id := range partUUIDs {
		if _, ok := qtyByPart[id]; !ok {
			unique = append(unique, id)
		}
		qtyByPart[id]++
	}

	parts, err := s.partRepo.GetByUUIDs(ctx, unique)
	if err != nil {
		metrics.OrdersCreateFailedTotal.WithLabelValues("repository_error").Inc()
		return nil, err
	}
	if len(parts) != len(unique) {
		metrics.OrdersCreateFailedTotal.WithLabelValues("part_not_found").Inc()
		return nil, model.ErrPartNotFound
	}

	partsByUUID := make(map[uuid.UUID]*model.Part, len(parts))
	for _, part := range parts {
		partsByUUID[part.UUID] = part
	}

	items := make([]model.OrderItem, 0, len(unique))
	var totalPrice float64
	expanded := make([]uuid.UUID, 0, len(partUUIDs))

	for _, id := range unique {
		part := partsByUUID[id]
		qty := qtyByPart[id]
		if part.InStock < qty {
			metrics.OrdersCreateFailedTotal.WithLabelValues("insufficient_stock").Inc()
			return nil, model.ErrInsufficientStock
		}
		items = append(items, model.OrderItem{
			PartUUID:  id,
			Quantity:  qty,
			UnitPrice: part.Price,
		})
		totalPrice += part.Price * float64(qty)
		for i := 0; i < qty; i++ {
			expanded = append(expanded, id)
		}
	}

	newOrder := &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUUIDs:  expanded,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	for i := range newOrder.Items {
		newOrder.Items[i].OrderUUID = newOrder.OrderUUID
	}

	if err := s.orderRepo.Create(ctx, newOrder); err != nil {
		metrics.OrdersCreateFailedTotal.WithLabelValues("repository_error").Inc()
		return nil, err
	}

	metrics.OrdersCreatedTotal.Inc()
	metrics.OrderTotalPrice.Observe(totalPrice)

	logger.Info(ctx, "order created",
		zap.String("order_uuid", newOrder.OrderUUID.String()),
		zap.String("user_uuid", userUUID.String()),
		zap.Float64("total_price", totalPrice),
	)

	return newOrder, nil
}
