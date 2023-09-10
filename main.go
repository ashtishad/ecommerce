package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/db/conn"
	"github.com/ashtishad/ecommerce/lib"
	productApp "github.com/ashtishad/ecommerce/product-api/cmd/app"
	usersApp "github.com/ashtishad/ecommerce/users-api/cmd/app"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// Initialize shared DB client
	l := initLogger()
	//sanityCheck(l)
	dbClient := initDBClient(l)
	defer dbClient.Close()

	// run db migrations if any
	m, err := migrate.New(
		"file://db/migrations",
		conn.GetDSNString(l),
	)
	if err != nil {
		l.Error("error creating migration: %v", "err", err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("error applying migration: %v", "err", err.Error())
	}

	//generate.GenerateUsers(dbClient, l, 1000)

	// Initialize first server
	userServer := initServerConfig("SERVER_PORT")
	wg.Add(1)
	go func() {
		usersApp.Start(userServer, dbClient, l)
		wg.Done()
	}()

	// Initialize second server
	productServer := initServerConfig("PRODUCT_API_PORT")
	wg.Add(1)
	go func() {
		productApp.Start(productServer, dbClient, l)
		wg.Done()
	}()

	// Wait for an interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Graceful shutdown for userServer
	log.Println("Shutting down user server...")
	if err := userServer.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the user server: %v\n", err)
	}
	log.Println("User server gracefully stopped")

	// Graceful shutdown for productServer
	log.Println("Shutting down product server...")
	if err := productServer.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the product server: %v\n", err)
	}
	log.Println("Product server gracefully stopped")

	// Wait for all goroutines to complete
	wg.Wait()
}

func initServerConfig(portEnv string) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv(portEnv)),
		IdleTimeout:    100 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func initLogger() *slog.Logger {
	handlerOpts := lib.GetSlogConf()
	l := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(l)
	return l
}

func initDBClient(l *slog.Logger) *sql.DB {
	conn := conn.GetDbClient(l)
	return conn
}
