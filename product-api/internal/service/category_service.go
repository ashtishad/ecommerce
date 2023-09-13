package service

import (
	"context"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
)

type CategoryService interface {
	NewCategory(ctx context.Context, req domain.NewCategoryRequestDTO) (*domain.CategoryResponseDTO, lib.APIError)
	NewSubCategory(ctx context.Context, req domain.NewCategoryRequestDTO, parentUUID string) (*domain.CategoryResponseDTO, lib.APIError)
	GetAllCategoriesByHierarchy(ctx context.Context) ([]*domain.CategoryResponseDTO, lib.APIError)
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

func (s *DefaultCategoryService) NewSubCategory(ctx context.Context, req domain.NewCategoryRequestDTO, parentUUID string) (*domain.CategoryResponseDTO, lib.APIError) {
	if apiErr := ValidateNewCategoryRequest(req); apiErr != nil {
		return nil, apiErr
	}

	subCategory := domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	createdSubCategory, apiErr := s.repo.CreateSubCategory(ctx, subCategory, parentUUID)
	if apiErr != nil {
		return nil, apiErr
	}

	response := createdSubCategory.ToCategoryResponseDTO()

	return response, nil
}

func (s *DefaultCategoryService) GetAllCategoriesByHierarchy(ctx context.Context) ([]*domain.CategoryResponseDTO, lib.APIError) {
	categories, apiErr := s.repo.GetAllCategoriesWithHierarchy(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	categoryResponseDTOs := make([]*domain.CategoryResponseDTO, 0, len(categories))
	for _, category := range categories {
		categoryResponseDTOs = append(categoryResponseDTOs, category.ToCategoryResponseDTO())
	}

	return categoryResponseDTOs, nil
}
