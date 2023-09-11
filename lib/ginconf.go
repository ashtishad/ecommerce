package lib

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logger implements logger util for gin
func Logger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] | %s | %s | %d | %s |%s\n",
		// param.ClientIP,
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		param.Method,
		param.Path,
		// param.Request.Proto,
		param.StatusCode,
		param.Latency,
		// param.Request.UserAgent(),
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
