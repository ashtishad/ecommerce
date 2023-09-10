package service

import (
	"context"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
)

type CategoryService interface {
	NewCategory(ctx context.Context, req domain.NewCategoryRequestDTO) (*domain.CategoryResponseDTO, lib.APIError)
}

type DefaultCategoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) *DefaultCategoryService {
	return &DefaultCategoryService{repo: repo}
}

func (s *DefaultCategoryService) NewCategory(ctx context.Context, req domain.NewCategoryRequestDTO) (*domain.CategoryResponseDTO, lib.APIError) {
	if apiErr := ValidateNewCategoryRequest(req); apiErr != nil {
		return nil, apiErr
	}

	category := domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	createdCategory, apiErr := s.repo.CreateCategory(ctx, category)
	if apiErr != nil {
		return nil, apiErr
	}

	response := createdCategory.ToCategoryResponseDTO()

	return response, nil
}
