package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/database"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func StartUsersAPI() {
	// init slog
	handlerOpts := lib.GetSlogConf()

	l := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(l)
	slog.Info("API NAME", "name", "users-api")

	// environment variables check
	sanityCheck(l)

	gin.SetMode(gin.ReleaseMode)
	var r = gin.New()

	// database connection config
	conn := database.GetDbClient(l)
	defer func(conn *sql.DB) {
		dbConnCloseErr := conn.Close()
		if dbConnCloseErr != nil {
			l.Error("error closing db connection", "err", dbConnCloseErr.Error())
			return
		}
	}(conn)

	// run db migrations if any
	m, err := migrate.New(
		"file://db/migrations",
		database.GetDSNString(l),
	)
	if err != nil {
		l.Error("error creating migration: %v", "err", err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("error applying migration: %v", "err", err.Error())
	}

	// database.GenerateUsers(conn, l, 1000)

	// wire up the handler
	userRepositoryDB := domain.NewUserRepositoryDB(conn, l)
	uh := UserHandlers{service.NewUserService(userRepositoryDB), l}

	// Server Config
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler:        r,
		IdleTimeout:    100 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// route url mappings
	setUsersApiRoutes(r, uh)

	// custom logger middleware
	r.Use(gin.LoggerWithFormatter(lib.Logger))

	// custom recovery middleware
	r.Use(gin.CustomRecovery(lib.Recover))

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("could not start server: %v\n", "err", err.Error())
		}
	}()

	// graceful shutdown
	lib.GracefulShutdown(srv)
}

func setUsersApiRoutes(r *gin.Engine, uh UserHandlers) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", uh.createUserHandler)
		userRoutes.PUT("/:user_id", uh.updateUserHandler)
		userRoutes.GET("", uh.GetUsersHandler)
	}
}
