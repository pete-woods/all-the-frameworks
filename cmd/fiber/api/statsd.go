package api

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gofiber/fiber/v2"
)

func statsdMiddleWare(stats *statsd.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		defer func() {
			tags := []string{
				fmt.Sprintf("request.method:%s", c.Method()),
				fmt.Sprintf("request.route:%s", c.Path()),
				fmt.Sprintf("response.status_code:%d", c.Response().StatusCode()),
			}
			_ = stats.Timing("handler", time.Since(start), tags, 1)
		}()
		return c.Next()
	}
}
