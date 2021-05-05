package api

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/DataDog/datadog-go/statsd"
	"golang.org/x/sync/errgroup"
)

type API struct {
	server *http.Server
}

func New(addr string, stats *statsd.Client) *API {
	r := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound:               http.HandlerFunc(notFound),
		MethodNotAllowed:       http.HandlerFunc(methodNotAllowed),
		PanicHandler:           panicHandler,
	}

	a := &API{
		server: &http.Server{
			Addr:         addr,
			Handler:      statsdMiddleWare(stats, logMiddleWare(r)),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	r.HandlerFunc("GET", "/v1/bananas", a.getBananas)
	r.HandlerFunc("GET", "/v1/banana/:id", a.getBanana)

	return a
}

func (a *API) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.server.Shutdown(sctx)
	})

	g.Go(a.server.ListenAndServe)

	return g.Wait()
}

type banana struct {
	ID string `json:"id"`
}

func (a *API) getBananas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	writeJSON(w, http.StatusOK, []banana{
		{
			ID: "abc",
		},
		{
			ID: "def",
		},
	})
}

func (a *API) getBanana(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "parameter 'id' is required")
		return
	}

	writeJSON(w, http.StatusOK, banana{
		ID: id,
	})
}
