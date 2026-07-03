package v1

import (
	"context"
	"inventory/internal/converter"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *api) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	part, err := h.partService.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Деталь с UUID %s не найдена", req.Uuid)
	}

	return &pb.GetPartResponse{
		Part: converter.ToPbPart(part),
	}, nil
}
