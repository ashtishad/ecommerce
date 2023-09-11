package app

import (
	"context"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
	"github.com/ashtishad/ecommerce/product-api/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type CategoryHandlers struct {
	service service.CategoryService
	l       *slog.Logger
}

func (ch *CategoryHandlers) CreateCategory(c *gin.Context) {
	var newCategoryReqDTO domain.NewCategoryRequestDTO
	if err := c.ShouldBindJSON(&newCategoryReqDTO); err != nil {
		ch.l.Error("failed to bind create category req dto", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), lib.TimeoutCreateCategory)
	defer cancel()

	createdCategory, apiErr := ch.service.NewCategory(timeoutCtx, newCategoryReqDTO)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	c.JSON(http.StatusOK, createdCategory)
}
