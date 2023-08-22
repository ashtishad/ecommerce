package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/domain"
	"github.com/ashtishad/ecommerce/pkg/ginconf"
	"github.com/ashtishad/ecommerce/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	sanityCheck()

	gin.SetMode(gin.ReleaseMode)
	var r = gin.New()

	// initiated logger, dependency injection, create once, inject it where needed
	l := log.New(os.Stdout, "users-api ", log.LstdFlags)

	// database connection config
	conn := getDbClient()
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			l.Printf("couldn't close the database client : %v", err.Error())
		}
	}(conn)

	// wire up the handler
	userRepositoryDB := domain.NewUserRepositoryDB(conn, l)
	uh := UserHandlers{service.NewUserService(userRepositoryDB), l}
	gh := GoogleAuthHandler{l: l}

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
	setRouteMappings(r, uh, gh)

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

func setRouteMappings(r *gin.Engine, uh UserHandlers, gh GoogleAuthHandler) {
	// Group routes related to users
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", uh.createUserHandler)
		userRoutes.PUT("/:user_id", uh.updateUserHandler)
		userRoutes.POST("/existing-user", uh.existingUserHandler)
		userRoutes.GET("/existing-user", uh.existingUserHandler)
	}

	// Group routes related to Google authentication
	// http://localhost:8000/google-auth/login
	googleAuthRoutes := r.Group("/google-auth")
	{
		googleAuthRoutes.GET("/login", gh.StartGoogleLoginHandler)
		googleAuthRoutes.POST("/callback", gh.GoogleCallbackHandler)
		googleAuthRoutes.GET("/callback", gh.GoogleCallbackHandler)
	}
}
