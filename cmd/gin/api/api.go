package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type API struct {
	server *http.Server
}

func New(addr string, stats *statsd.Client) *API {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		statsdMiddleWare(stats),
	)

	a := &API{
		server: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	r.UseRawPath = true
	r.GET("/v1/bananas", a.getBananas)
	r.GET("/v1/banana/:id", a.getBanana)

	authorized := r.Group("", gin.BasicAuth(gin.Accounts{
		"foo":  "bar",
		"manu": "123",
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
		return a.server.Shutdown(cctx)
	})

	g.Go(func() error {
		err := a.server.ListenAndServe()
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

func (a *API) getBananas(c *gin.Context) {
	c.JSON(http.StatusOK, []banana{
		{
			ID: "abc",
		},
		{
			ID: "def",
		},
	})
}

func (a *API) getBanana(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, banana{
		ID: id,
	})
}

func (a *API) postAdmin(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	var v struct {
		Value string `json:"value" binding:"required"`
	}

	if c.Bind(&v) == nil {
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"value": v.Value,
		})
	}
}
