package v1

import (
	"context"
	"errors"

	"order/internal/converter"
	"order/internal/model"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) ListParts(ctx context.Context) (orderV1.ListPartsRes, error) {
	parts, err := h.partService.ListParts(ctx)
	if err != nil {
		return nil, err
	}
	result := make(orderV1.ListPartsOKApplicationJSON, 0, len(parts))
	for _, part := range parts {
		result = append(result, *converter.PartToDto(part))
	}
	return &result, nil
}

func (h *api) CreatePart(ctx context.Context, req *orderV1.CreatePartRequest) (orderV1.CreatePartRes, error) {
	category := req.Category.Or("GENERAL")
	inStock := req.InStock.Or(0)

	part, err := h.partService.CreatePart(ctx, req.Name, req.Price, category, inStock)
	if err != nil {
		if errors.Is(err, model.ErrPartAlreadyExists) {
			return &orderV1.CreatePartBadRequest{Code: 400, Message: err.Error()}, nil
		}
		return &orderV1.CreatePartInternalServerError{Code: 500, Message: err.Error()}, nil
	}
	return converter.PartToDto(part), nil
}

func (h *api) GetPartByUUID(ctx context.Context, params orderV1.GetPartByUUIDParams) (orderV1.GetPartByUUIDRes, error) {
	part, err := h.partService.GetPartByUUID(ctx, params.PartUUID)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.GetPartByUUIDNotFound{Code: 404, Message: "Part not found"}, nil
		}
		return nil, err
	}
	return converter.PartToDto(part), nil
}

func (h *api) UpdatePart(ctx context.Context, req *orderV1.UpdatePartRequest, params orderV1.UpdatePartParams) (orderV1.UpdatePartRes, error) {
	part, err := h.partService.UpdatePart(ctx, params.PartUUID, req.Name, req.Price, req.Category, req.InStock)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.UpdatePartNotFound{Code: 404, Message: "Part not found"}, nil
		}
		return &orderV1.UpdatePartInternalServerError{Code: 500, Message: err.Error()}, nil
	}
	return converter.PartToDto(part), nil
}

func (h *api) DeletePart(ctx context.Context, params orderV1.DeletePartParams) (orderV1.DeletePartRes, error) {
	err := h.partService.DeletePart(ctx, params.PartUUID)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.DeletePartNotFound{Code: 404, Message: "Part not found"}, nil
		}
		return &orderV1.DeletePartInternalServerError{Code: 500, Message: err.Error()}, nil
	}
	return &orderV1.DeletePartNoContent{}, nil
}

func (h *api) ListPaymentsByOrder(ctx context.Context, params orderV1.ListPaymentsByOrderParams) (orderV1.ListPaymentsByOrderRes, error) {
	payments, err := h.orderService.ListPaymentsByOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.ListPaymentsByOrderNotFound{Code: 404, Message: "Order not found"}, nil
		}
		return nil, err
	}
	result := make(orderV1.ListPaymentsByOrderOKApplicationJSON, 0, len(payments))
	for _, payment := range payments {
		result = append(result, *converter.PaymentToDto(payment))
	}
	return &result, nil
}

func (h *api) GetPaymentByUUID(ctx context.Context, params orderV1.GetPaymentByUUIDParams) (orderV1.GetPaymentByUUIDRes, error) {
	payment, err := h.orderService.GetPaymentByUUID(ctx, params.PaymentUUID)
	if err != nil {
		if errors.Is(err, model.ErrPaymentNotFound) {
			return &orderV1.GetPaymentByUUIDNotFound{Code: 404, Message: "Payment not found"}, nil
		}
		return nil, err
	}
	return converter.PaymentToDto(payment), nil
}
