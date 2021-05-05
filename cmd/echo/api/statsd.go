package api

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/DataDog/datadog-go/statsd"
)

func statsdMiddleWare(stats *statsd.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			start := time.Now()
			defer func() {
				tags := []string{
					fmt.Sprintf("request.method:%s", c.Request().Method),
					fmt.Sprintf("request.route:%s", c.Path()),
					fmt.Sprintf("response.status_code:%d", c.Response().Status),
				}
				_ = stats.Timing("handler", time.Since(start), tags, 1)
				_ = stats.Incr("handler_count", tags, 1)
			}()
			return next(c)
		}
	}
}
