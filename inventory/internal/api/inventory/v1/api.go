package v1

import (
	"inventory/internal/service"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"
)

type api struct {
	pb.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewAPI(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
