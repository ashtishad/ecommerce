package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/users-api/database"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/ashtishad/ecommerce/users-api/pkg/ginconf"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func StartUsersAPI() {
	// initiated logger, dependency injection, create once, inject it where needed
	l := log.New(os.Stdout, "users-api ", log.LstdFlags)

	sanityCheck(l)

	gin.SetMode(gin.ReleaseMode)
	var r = gin.New()

	// database connection config
	conn := database.GetMySqlDBClient()
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			l.Printf("couldn't close the database client : %v", err.Error())
		}
	}(conn)

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
	r.Use(gin.LoggerWithFormatter(ginconf.Logger))

	// custom recovery middleware
	r.Use(gin.CustomRecovery(ginconf.Recover))

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatalf("could not start server: %v\n", err)
		}
	}()

	// graceful shutdown
	ginconf.GracefulShutdown(srv)
}

func setUsersApiRoutes(r *gin.Engine, uh UserHandlers) {
	// Group routes related to users
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", uh.createUserHandler)
		userRoutes.PUT("/:user_id", uh.updateUserHandler)
		userRoutes.POST("/existing-user", uh.existingUserHandler)
		userRoutes.GET("/existing-user", uh.existingUserHandler)
	}
}
