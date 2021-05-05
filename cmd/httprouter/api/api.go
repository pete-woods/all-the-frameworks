package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/julienschmidt/httprouter"
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

	r.HandlerFunc("POST", "/v1/admin", basicAuth(a.postAdmin, map[string]string{
		"foo":  "bar",
		"manu": "123",
	}, "realm"))

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

func (a *API) postAdmin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := r.Context().Value("username").(string)

	var v struct {
		Value string `json:"value" binding:"required"`
	}

	d := json.NewDecoder(r.Body)
	err := d.Decode(&v)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"value": v.Value,
	})
}
