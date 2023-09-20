package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ashtishad/ecommerce/product-api/internal/service"
	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	s service.DefaultBrandService
}

func (b *BrandHandler) GetAllBrands(c *gin.Context) {
	status := c.DefaultQuery("status", "active")

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
	defer cancel()

	brands, tc, err := b.s.GetAll(timeoutCtx, status)
	if err != nil {
		c.JSON(err.StatusCode(), err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"brands":     brands,
		"totalCount": tc,
	})
}
