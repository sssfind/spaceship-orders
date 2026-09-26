package v1

import "order/internal/service"

type api struct {
	orderService service.OrderService
	partService  service.PartService
}

func NewAPI(orderService service.OrderService, partService service.PartService) *api {
	return &api{
		orderService: orderService,
		partService:  partService,
	}
}
