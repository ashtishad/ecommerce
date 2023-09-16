package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ashtishad/ecommerce/lib"
	productApp "github.com/ashtishad/ecommerce/product-api/cmd/app"
	usersApp "github.com/ashtishad/ecommerce/users-api/cmd/app"
)

func main() {
	var wg sync.WaitGroup

	l := lib.InitSlogger()

	lib.SanityCheck(l)

	dbClient := lib.InitDBClient(l)
	defer func(dbClient *sql.DB) {
		if dbClsErr := dbClient.Close(); dbClsErr != nil {
			l.Error("unable to close db", "err", dbClsErr)
			os.Exit(1)
		}
	}(dbClient)

	// generate.Users(dbClient, l, 1000)

	userServer := lib.InitServerConfig("USER_API_PORT")

	wg.Add(1)

	go func() {
		usersApp.Start(userServer, dbClient, l)
		wg.Done()
	}()

	productServer := lib.InitServerConfig("PRODUCT_API_PORT")

	wg.Add(1)

	go func() {
		productApp.Start(productServer, dbClient, l)
		wg.Done()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	wg.Add(2)

	go lib.GracefulShutdown(ctx, userServer, &wg, "User")
	go lib.GracefulShutdown(ctx, productServer, &wg, "Product")

	wg.Wait()
}
