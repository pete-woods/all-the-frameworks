package api

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
)

func statsdMiddleWare(stats *statsd.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			tags := []string{
				fmt.Sprintf("request.method:%s", c.Request.Method),
				fmt.Sprintf("request.route:%s", c.FullPath()),
				fmt.Sprintf("response.status_code:%d", c.Writer.Status()),
			}
			_ = stats.Timing("handler", time.Since(start), tags, 1)
			_ = stats.Incr("handler_count", tags, 1)
		}()
		c.Next()
	}
}
