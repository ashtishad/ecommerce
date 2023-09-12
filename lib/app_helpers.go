package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ashtishad/ecommerce/db/conn"
	"github.com/golang-migrate/migrate/v4"
)

func InitServerConfig(portEnv string) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv(portEnv)),
		IdleTimeout:    100 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func InitSlogger() *slog.Logger {
	handlerOpts := getSlogConf()
	l := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(l)

	return l
}

func InitDBClient(l *slog.Logger) *sql.DB {
	dbClient := conn.GetDBClient(l)

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

	return dbClient
}

func GracefulShutdown(ctx context.Context, srv *http.Server, wg *sync.WaitGroup, serverName string) {
	defer wg.Done()
	log.Printf("Shutting down %s server...\n", serverName)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the %s server: %v\n", serverName, err)
	}

	log.Printf("%s server gracefully stopped\n", serverName)
}

// SanityCheck checks that all required environment variables are set.
// if any of the required variables is not defined, it prints a log message.
func SanityCheck(l *slog.Logger) {
	envProps := []string{
		"SERVER_ADDRESS",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
		"PRODUCT_API_PORT",
		"USER_API_PORT",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			l.Warn(fmt.Sprintf("environment variable %s not defined", k))
		}
	}
}
