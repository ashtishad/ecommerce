package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
)

type BrandService interface {
	GetAll(ctx context.Context, status string) ([]domain.BrandResponseDTO, int, lib.APIError)
}

type DefaultBrandService struct {
	repo domain.BrandRepository
}

func NewBrandService(repo domain.BrandRepository) DefaultBrandService {
	return DefaultBrandService{repo: repo}
}

func (s *DefaultBrandService) GetAll(ctx context.Context, status string) ([]domain.BrandResponseDTO, int, lib.APIError) {
	const statusRegex = `^(active|inactive|deleted)$`
	if m := regexp.MustCompile(statusRegex).MatchString(strings.ToLower(status)); !m {
		return nil, 0, lib.NewBadRequestError("brand status should be active, inactive or deleted")
	}

	brands, err := s.repo.GetAllBrands(ctx, status)
	if err != nil {
		return nil, 0, err
	}

	totalCount := len(brands)

	res := make([]domain.BrandResponseDTO, 0, totalCount)
	for _, brand := range brands {
		res = append(res, *brand.ToBrandResponseDTO())
	}

	return res, totalCount, nil
}
