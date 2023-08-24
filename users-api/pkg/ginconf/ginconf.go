package ginconf

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// GracefulShutdown wait for interrupt signal to gracefully shut down the server with a timeout of 1 Minute.
func GracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// graceful shutdown
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("Server gracefully stopped")
}

// Logger implements logger util for gin
func Logger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] | %s | %s | %d | %s |%s\n",
		//param.ClientIP,
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		param.Method,
		param.Path,
		//param.Request.Proto,
		param.StatusCode,
		param.Latency,
		//param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// Recover from any panics and writes a 500 if there was one.
func Recover(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
