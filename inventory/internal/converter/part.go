package converter

import (
	"inventory/internal/model"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"
)

// ToDomainFilter переводит Protobuf-фильтр в чистую доменную модель фильтра
func ToDomainFilter(pbFilter *pb.PartsFilter) *model.PartsFilter {
	if pbFilter == nil {
		return nil
	}

	filter := &model.PartsFilter{
		UUIDs:                 pbFilter.Uuids,
		Names:                 pbFilter.Names,
		ManufacturerCountries: pbFilter.ManufacturerCountries,
		Tags:                  pbFilter.Tags,
	}

	if len(pbFilter.Categories) > 0 {
		filter.Categories = make([]model.Category, 0, len(pbFilter.Categories))
		for _, cat := range pbFilter.Categories {
			filter.Categories = append(filter.Categories, model.Category(cat))
		}
	}

	return filter
}

// ToPbPart переводит одну чистую доменную деталь в структуру Protobuf для сети
func ToPbPart(p *model.Part) *pb.Part {
	if p == nil {
		return nil
	}

	return &pb.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Description:   p.Description,
		Category:      pb.Category(p.Category),
		Manufacturer: &pb.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Dimensions: &pb.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
	}
}

// ToPbParts конвертирует слайс доменных моделей в слайс структур Protobuf
func ToPbParts(parts []model.Part) []*pb.Part {
	if parts == nil {
		return nil
	}

	pbParts := make([]*pb.Part, 0, len(parts))
	for _, p := range parts {
		pbParts = append(pbParts, ToPbPart(&p))
	}

	return pbParts
}
