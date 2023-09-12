package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Start(srv *http.Server, dbClient *sql.DB, l *slog.Logger) {
	gin.SetMode(gin.ReleaseMode)

	var r = gin.New()
	srv.Handler = r

	// wire up the handler
	userRepositoryDB := domain.NewUserRepositoryDB(dbClient, l)
	uh := UserHandlers{service.NewUserService(userRepositoryDB), l}

	// route url mappings
	setUsersAPIRoutes(r, uh)

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

func setUsersAPIRoutes(r *gin.Engine, uh UserHandlers) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", uh.createUserHandler)
		userRoutes.PUT("/:user_id", uh.updateUserHandler)
		userRoutes.GET("", uh.GetUsersHandler)
	}
}
