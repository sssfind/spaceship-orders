package part

import (
	"context"

	domainModel "inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (r *repo) List(ctx context.Context, filter *domainModel.PartsFilter) ([]domainModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filterTagsSet map[string]struct{}
	if filter != nil && len(filter.Tags) > 0 {
		filterTagsSet = make(map[string]struct{}, len(filter.Tags))
		for _, tag := range filter.Tags {
			filterTagsSet[tag] = struct{}{}
		}
	}

	var result []domainModel.Part

	for _, part := range r.parts {
		if filter == nil {
			result = append(result, converter.PartToDomain(part))
			continue
		}

		if len(filter.UUIDs) > 0 && !containsString(filter.UUIDs, part.UUID) {
			continue
		}
		if len(filter.Names) > 0 && !containsString(filter.Names, part.Name) {
			continue
		}

		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
			continue
		}

		if len(filter.ManufacturerCountries) > 0 && !containsString(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}

		if len(filter.Tags) > 0 && !containsTags(filter.Tags, filterTagsSet) {
			continue
		}

		result = append(result, converter.PartToDomain(part))
	}
	return result, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func containsCategory(categories []domainModel.Category, target repoModel.Category) bool {
	for _, item := range categories {
		if repoModel.Category(item) == target {
			return true
		}
	}
	return false
}

func containsTags(tags []string, target map[string]struct{}) bool {
	for _, tag := range tags {
		if _, exists := target[tag]; exists {
			return true
		}
	}
	return false
}
