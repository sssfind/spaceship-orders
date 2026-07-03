package v1

import "order/internal/service"

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
