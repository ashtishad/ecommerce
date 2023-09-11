package main

import (
	"context"
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
	defer dbClient.Close()

	// generate.GenerateUsers(dbClient, l, 1000)

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

	go lib.GracefulShutdown(userServer, ctx, &wg, "User")
	go lib.GracefulShutdown(productServer, ctx, &wg, "Product")

	wg.Wait()
}
