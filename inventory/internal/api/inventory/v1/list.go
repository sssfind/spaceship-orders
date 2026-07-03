package v1

import (
	"context"

	"inventory/internal/converter"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"
)

func (h *api) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	domainFilter := converter.ToDomainFilter(req.Filter)

	parts, err := h.partService.ListParts(ctx, domainFilter)
	if err != nil {
		return nil, err
	}

	return &pb.ListPartsResponse{
		Parts: converter.ToPbParts(parts),
	}, nil
}
