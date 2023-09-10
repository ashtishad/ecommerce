package app

import (
	"context"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
	"github.com/ashtishad/ecommerce/product-api/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
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

	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	createdCategory, apiErr := ch.service.NewCategory(ctx, newCategoryReqDTO)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	c.JSON(http.StatusOK, createdCategory)
}
