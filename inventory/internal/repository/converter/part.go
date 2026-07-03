package converter

import (
	domainModel "inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)

func PartToDomain(part *repoModel.Part) domainModel.Part {
	if part == nil {
		return domainModel.Part{}
	}
	return domainModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Description:   part.Description,
		Category:      domainModel.Category(part.Category),
		Dimensions:    domainModel.Dimensions(part.Dimensions),
		Manufacturer:  domainModel.Manufacturer(part.Manufacturer),
		Tags:          []string{},
	}
}
