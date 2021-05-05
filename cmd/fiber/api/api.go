package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type API struct {
	app  *fiber.App
	addr string
}

func New(addr string, stats *statsd.Client) *API {
	gin.SetMode(gin.ReleaseMode)
	a := &API{
		app: fiber.New(fiber.Config{
			ReadTimeout:           10 * time.Second,
			WriteTimeout:          10 * time.Second,
			DisableStartupMessage: true,
			ErrorHandler:          jsonErrorHandler,
		}),
		addr: addr,
	}

	a.app.Use(logger.New())
	a.app.Use(statsdMiddleWare(stats))

	a.app.Get("/v1/bananas", a.getBananas)
	a.app.Get("/v1/banana/:id", a.getBanana)

	return a
}

var jsonErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	type httpError struct {
		Error string `json:"error"`
	}
	return c.Status(code).JSON(httpError{Error: err.Error()})
}

func (a *API) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		return a.app.Shutdown()
	})

	g.Go(func() error {
		return a.app.Listen(a.addr)
	})

	return g.Wait()
}

type banana struct {
	ID string `json:"id"`
}

func (a *API) getBananas(c *fiber.Ctx) error {
	return c.JSON([]banana{
		{
			ID: "abc",
		},
		{
			ID: "def",
		},
	})
}

func (a *API) getBanana(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(banana{
		ID: id,
	})
}
