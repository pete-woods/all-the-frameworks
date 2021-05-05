package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/DataDog/datadog-go/statsd"
	"golang.org/x/sync/errgroup"
)

type API struct {
	echo *echo.Echo
	addr string
}

func New(addr string, stats *statsd.Client) *API {
	e := echo.New()

	e.Server = &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	e.TLSServer = e.Server
	e.HideBanner = true
	e.Use(statsdMiddleWare(stats))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	a := &API{
		echo: e,
		addr: addr,
	}

	e.GET("/v1/bananas", a.getBananas)
	e.GET("/v1/banana/:id", a.getBanana)

	accounts := map[string]string{
		"foo":  "bar",
		"manu": "123",
	}
	authorized := e.Group("", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if accounts[username] == password {
			c.Set("username", username)
			return true, nil
		}
		return false, nil
	}))

	authorized.POST("/v1/admin", a.postAdmin)

	return a
}

func (a *API) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.echo.Shutdown(cctx)
	})

	g.Go(func() error {
		err := a.echo.Start(a.addr)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return g.Wait()
}

type banana struct {
	ID string `json:"id"`
}

func (a *API) getBananas(c echo.Context) error {
	return c.JSON(http.StatusOK, []banana{
		{
			ID: "abc",
		},
		{
			ID: "def",
		},
	})
}

func (a *API) getBanana(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, banana{
		ID: id,
	})
}

func (a *API) postAdmin(c echo.Context) error {
	user := c.Get("username").(string)

	var v struct {
		Value string `json:"value" binding:"required"`
	}

	err := c.Bind(&v)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user":  user,
		"value": v.Value,
	})
}
