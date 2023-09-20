package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
	"github.com/ashtishad/ecommerce/product-api/internal/service"
	"github.com/gin-gonic/gin"
)

func Start(srv *http.Server, dbClient *sql.DB, l *slog.Logger) {
	gin.SetMode(gin.ReleaseMode)

	var r = gin.New()
	srv.Handler = r

	// wire up the handlers
	categoryRepoDB := domain.NewCategoryRepoDB(dbClient, l)
	ch := CategoryHandlers{
		service: service.NewCategoryService(categoryRepoDB),
		l:       l,
	}

	brandRepoDB := domain.NewBrandRepoDB(dbClient, l)
	bh := BrandHandler{s: service.NewBrandService(brandRepoDB)}

	// route url mappings
	setProductAPIRoutes(r, ch)
	setBrandAPIRoutes(r, bh)
	// custom logger middleware
	r.Use(gin.LoggerWithFormatter(lib.Logger))

	// custom recovery middleware
	r.Use(gin.CustomRecovery(lib.Recover))

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("could not start server: %v\n", "err", err.Error(), "srv", srv.Addr)
		}
	}()
}

func setProductAPIRoutes(r *gin.Engine, ch CategoryHandlers) {
	categoriesRoutes := r.Group("/categories")
	{
		categoriesRoutes.GET("", ch.GetAllCategories)
		categoriesRoutes.POST("", ch.CreateCategory)
		categoriesRoutes.POST("/:category_id/subcategories", ch.CreateSubCategory)
	}
}

func setBrandAPIRoutes(r *gin.Engine, bh BrandHandler) {
	r.GET("/brands", bh.GetAllBrands)
}
